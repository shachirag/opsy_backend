package logEntry

import (
	"fmt"
	"math"
	"opsy_backend/database"
	"opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch daily insights data for a specified month
// @Description Retrieves mental health and physical health data for each date within a specified month
// @Tags logEntry
// @Accept json
// @Produce json
// @Param userId query string true "user ID"
// @Param monthYear query string true "Month and Year (MM-YYYY)"
// @Success 200 {object} logEntry.MonthlyInsightsResDto
// @Router /user/monthly-insights [get]
// Wrapper function for DailyInsights with onlyEvenDates parameter
func MonthlyInsights(c *fiber.Ctx) error {
	onlyEvenDates := true

	return DailyInsights(c, onlyEvenDates)
}

func DailyInsights(c *fiber.Ctx, onlyEvenDates bool) error {
	monthYear := c.Query("monthYear")
	parsedMonthYear, err := time.Parse("01-2006", monthYear)
	if err != nil {
		return fmt.Errorf("Invalid monthYear format: %s", err.Error())
	}

	startDate := time.Date(parsedMonthYear.Year(), parsedMonthYear.Month(), 2, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	mentalHealthData, err := MentalHealthInsightDaysData(c, startDate, endDate, onlyEvenDates)
	if err != nil {
		return err
	}
	fmt.Printf("Mental Health Data: %v\n", mentalHealthData)

	physicalHealthData, err := PhysicalHealthInsightDaysData(c, startDate, endDate, onlyEvenDates)
	if err != nil {
		return err
	}
	fmt.Printf("Physical Health Data: %v\n", physicalHealthData)

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

func MentalHealthInsightDaysData(c *fiber.Ctx, startDate, endDate time.Time, onlyEvenDates bool) ([]logEntry.MonthlyMentalHealthRes, error) {
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
		"type":      "Mental Health",
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

	dailyData := make(map[string]logEntry.MonthlyMentalHealthRes)

	for _, entry := range logEntryData {
		dateKey := entry.When.Format("2006-01-02")

		if onlyEvenDates {
			day := entry.When.Day()
			if day%2 != 0 {
				continue
			}
		}

		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = logEntry.MonthlyMentalHealthRes{
				Date:    dateKey,
				AvgFeel: 0,
			}
		}

		temp := dailyData[dateKey]
		temp.AvgFeel += float64(mapFeel[entry.Feel])
		dailyData[dateKey] = temp
	}

	for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 2) {
		dateKey := day.Format("2006-01-02")
		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = logEntry.MonthlyMentalHealthRes{
				Date:    dateKey,
				AvgFeel: 0,
			}
		}
	}

	var dailyMentalHealthData []logEntry.MonthlyMentalHealthRes

	for _, data := range dailyData {
		total := data.AvgFeel
		count := len(logEntryData)

		avg := total / float64(count)
		avg = math.Round(avg*100) / 100
		dailyMentalHealthData = append(dailyMentalHealthData, logEntry.MonthlyMentalHealthRes{
			Date:    data.Date,
			AvgFeel: math.Round(math.Abs(avg)*100) / 100,
		})
	}

	sort.Slice(dailyMentalHealthData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", dailyMentalHealthData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", dailyMentalHealthData[j].Date)
		return dateI.Before(dateJ)
	})

	return dailyMentalHealthData, nil
}

func PhysicalHealthInsightDaysData(c *fiber.Ctx, startDate, endDate time.Time, onlyEvenDates bool) ([]logEntry.MonthlyPhysicalHealthRes, error) {
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
		"type":      "Physical Health",
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

	dailyData := make(map[string]logEntry.MonthlyPhysicalHealthRes)

	for _, entry := range logEntryData {
		dateKey := entry.When.Format("2006-01-02")

		if onlyEvenDates {
			day := entry.When.Day()
			if day%2 != 0 {
				continue
			}
		}

		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = logEntry.MonthlyPhysicalHealthRes{
				Date:         dateKey,
				AvgPainLevel: 0.0,
			}
		}

		temp := dailyData[dateKey]
		temp.AvgPainLevel += float64(entry.PainLevel)
		dailyData[dateKey] = temp
	}

	for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 2) {
		dateKey := day.Format("2006-01-02")
		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = logEntry.MonthlyPhysicalHealthRes{
				Date:         dateKey,
				AvgPainLevel: 0.0,
			}
		}
	}

	var dailyPhysicalHealthData []logEntry.MonthlyPhysicalHealthRes

	for _, data := range dailyData {
		total := data.AvgPainLevel
		count := len(logEntryData)

		avg := total / float64(count)
		avg = math.Round(avg*100) / 100
		dailyPhysicalHealthData = append(dailyPhysicalHealthData, logEntry.MonthlyPhysicalHealthRes{
			Date:         data.Date,
			AvgPainLevel: avg,
		})
	}

	sort.Slice(dailyPhysicalHealthData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", dailyPhysicalHealthData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", dailyPhysicalHealthData[j].Date)
		return dateI.Before(dateJ)
	})

	return dailyPhysicalHealthData, nil
}
