package logEntry

import (
	"fmt"
	"math"
	"opsy_backend/database"
	"opsy_backend/dto/users/logEntry"
	"time"

	"opsy_backend/entity"
	"sort"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch monthly insights data
// @Description Retrieves mental health and physical health data for a specified month
// @Tags logEntry
// @Accept json
// @Produce json
// @Param monthYear query string true "Month and Year (YYYY-MM)"
// @Success 200 {object} logEntry.InsightsResDto
// @Router /user/monthly-insights [get]
func MonthlyInsights(c *fiber.Ctx) error {
	mentalHealthData, err := MentalHealthInsightMonthsData(c)
	if err != nil {
		return err
	}

	physicalHealthData, err := PhysicalHealthInsightMonthsData(c)
	if err != nil {
		return err
	}

	response := logEntry.MonthlyInsightsResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: logEntry.MonthlyInsightsData{
			MentalHealth:   mentalHealthData,
			PhysicalHealth: physicalHealthData,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func MentalHealthInsightMonthsData(c *fiber.Ctx) ([]logEntry.MonthlyMentalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	monthYear := c.Query("monthYear")
	parsedMonthYear, err := time.Parse("01-2006", monthYear)
	if err != nil {
		return nil, fmt.Errorf("Invalid monthYear format: %s", err.Error())
	}

	startDate := time.Date(parsedMonthYear.Year(), parsedMonthYear.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endDate},
	}
	sortOptions := options.Find().SetSort(bson.M{"when": 1})

	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch data: %s", err.Error())
	}
	defer cursor.Close(ctx)

	var logEntryData []entity.LogEntryEntity
	if err := cursor.All(ctx, &logEntryData); err != nil {
		return nil, fmt.Errorf("Failed to decode data: %s", err.Error())
	}

	dailyData := make(map[string][]float64)
	dayCount := make(map[string]int)

	// Process log entry data
	for _, entry := range logEntryData {
		dateKey := entry.When.Format("2006-01-02")
		dayIndex := int(entry.When.Weekday())

		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = make([]float64, 7)
		}

		dailyData[dateKey][dayIndex] += float64(mapFeel[entry.Feel])
		dayCount[dateKey]++
	}

	var monthlyData []logEntry.MonthlyMentalHealthRes

	for date, values := range dailyData {
		dayIndex := int(startDate.AddDate(0, 0, 0).Weekday())
		dateTime, _ := time.Parse("2006-01-02", date)
		avg := values[dayIndex] / float64(dayCount[date])

		monthlyData = append(monthlyData, logEntry.MonthlyMentalHealthRes{
			Date:    dateTime.Format("2006-01-02"),
			AvgFeel: math.Round(math.Abs(avg)*100) / 100,
		})
	}

	sort.Slice(monthlyData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", monthlyData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", monthlyData[j].Date)
		return dateI.Before(dateJ)
	})

	for i := range monthlyData {
		if monthlyData[i].AvgFeel == 0 {
			monthlyData[i].AvgFeel = math.Abs(monthlyData[i].AvgFeel)
		} else if dayCount[monthlyData[i].Date] != 0 {
			monthlyData[i].AvgFeel *= -1
		}
	}

	return monthlyData, nil
}


func PhysicalHealthInsightMonthsData(c *fiber.Ctx) ([]logEntry.MonthlyPhysicalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	monthYear := c.Query("monthYear")
	parsedMonthYear, err := time.Parse("01-2006", monthYear)
	if err != nil {
		return nil, fmt.Errorf("Invalid monthYear format: %s", err.Error())
	}

	startDate := time.Date(parsedMonthYear.Year(), parsedMonthYear.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endDate},
	}

	sortOptions := options.Find().SetSort(bson.M{"when": 1})

	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch data: %s", err.Error())
	}
	defer cursor.Close(ctx)

	var logEntryData []entity.LogEntryEntity
	if err := cursor.All(ctx, &logEntryData); err != nil {
		return nil, fmt.Errorf("Failed to decode data: %s", err.Error())
	}

	dailyData := make(map[string][]float64)
	dayCount := make(map[string]int)

	// Process log entry data
	for _, entry := range logEntryData {
		dateKey := entry.When.Format("2006-01-02")
		dayIndex := int(entry.When.Weekday())

		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = make([]float64, 7)
		}

		// Using an array to keep track of values for each day
		dailyData[dateKey][dayIndex] += float64(entry.PainLevel)
		dayCount[dateKey]++
	}

	var monthlyData []logEntry.MonthlyPhysicalHealthRes

	// Calculate average pain level for each day in the month
	for date, values := range dailyData {
		dayIndex := int(startDate.AddDate(0, 0, 0).Weekday())
		dateTime, _ := time.Parse("2006-01-02", date)
		avg := values[dayIndex] / float64(dayCount[date])

		monthlyData = append(monthlyData, logEntry.MonthlyPhysicalHealthRes{
			Date:         dateTime.Format("2006-01-02"),
			AvgPainLevel: &avg,
		})
	}

	// Sort the data by date
	sort.Slice(monthlyData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", monthlyData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", monthlyData[j].Date)
		return dateI.Before(dateJ)
	})

	return monthlyData, nil
}
