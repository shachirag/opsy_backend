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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch monthly insights data
// @Description Retrieves mental health and physical health data for a specified month
// @Tags logEntry
// @Accept json
// @Produce json
// @Param userId query string true "user ID"
// @Param monthYear query string true "Month and Year (MM-YYYY)"
// @Success 200 {object} logEntry.InsightsResDto
// @Router /user/monthly-insights [get]
func MonthlyInsights(c *fiber.Ctx) error {

	monthYear := c.Query("monthYear")
	parsedMonthYear, err := time.Parse("01-2006", monthYear)
	if err != nil {
		return fmt.Errorf("Invalid monthYear format: %s", err.Error())
	}

	startDate := time.Date(parsedMonthYear.Year(), parsedMonthYear.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	mentalHealthData, err := MentalHealthInsightMonthsData(c, startDate, endDate)
	if err != nil {
		return err
	}

	physicalHealthData, err := PhysicalHealthInsightMonthsData(c, startDate, endDate)
	if err != nil {
		return err
	}

	allDates := make(map[string]bool)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		allDates[d.Format("2006-01-02")] = true
	}

	// Fill in missing data with 0 values for mental health
	for _, data := range mentalHealthData {
		delete(allDates, data.Date)
	}
	for date := range allDates {
		mentalHealthData = append(mentalHealthData, logEntry.MonthlyMentalHealthRes{
			Date:    date,
			AvgFeel: 0,
		})
	}

	// Fill in missing data with 0 values for physical health
	for _, data := range physicalHealthData {
		delete(allDates, data.Date)
	}
	for date := range allDates {
		physicalHealthData = append(physicalHealthData, logEntry.MonthlyPhysicalHealthRes{
			Date:         date,
			AvgPainLevel: new(float64), // Initializing with 0 value pointer
		})
	}

	sort.Slice(mentalHealthData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", mentalHealthData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", mentalHealthData[j].Date)
		return dateI.Before(dateJ)
	})

	sort.Slice(physicalHealthData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", physicalHealthData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", physicalHealthData[j].Date)
		return dateI.Before(dateJ)
	})

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

func MentalHealthInsightMonthsData(c *fiber.Ctx, startDate, endDate time.Time) ([]logEntry.MonthlyMentalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	userId := c.Query("userId")
	userObjID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, fmt.Errorf("invalid object id: %s", err.Error())
	}

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endDate},
		"userId":    userObjID,
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

	return monthlyData, nil
}

func PhysicalHealthInsightMonthsData(c *fiber.Ctx, startDate, endDate time.Time) ([]logEntry.MonthlyPhysicalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	userId := c.Query("userId")
	userObjID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, fmt.Errorf("invalid object id: %s", err.Error())
	}

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endDate},
		"userId":    userObjID,
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

	return monthlyData, nil
}
