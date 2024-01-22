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

// @Summary Fetch weekly insights data
// @Description Retrieves mental health and physical health data for a specified week
// @Tags logEntry
// @Accept json
// @Produce json
// @Param userId query string true "user ID"
// @Param firstDate query string true "Start date of the week (YYYY-MM-DD)"
// @Param endDate query string true "End date of the week (YYYY-MM-DD)"
// @Success 200 {object} logEntry.InsightsResDto
// @Router /user/weekly-insights [get]
func WeeklyInsights(c *fiber.Ctx) error {
	mentalHealthData, err := MentalHealthInsightWeeksData(c)
	if err != nil {
		return err
	}

	physicalHealthData, err := PhysicalHealthInsightWeeksData(c)
	if err != nil {
		return err
	}

	response := logEntry.InsightsResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: logEntry.InsightsData{
			MentalHealth:   mentalHealthData,
			PhysicalHealth: physicalHealthData,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

var mapFeel = map[string]int{
	"In Crisis":  1,
	"Struggling": 2,
	"Surviving":  3,
	"Thriving":   4,
	"Excelling":  5,
}

func MentalHealthInsightWeeksData(c *fiber.Ctx) ([]logEntry.MentalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	weekStartDate := c.Query("firstDate")
	weekEndDate := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("Invalid startDate format: %s", err.Error())
	}
	endDate, err := time.Parse("2006-01-02", weekEndDate)
	if err != nil {
		return nil, fmt.Errorf("Invalid endDate format: %s", err.Error())
	}
	endOfDay := endDate.Add(24 * time.Hour)

	userId := c.Query("userId")
	userObjID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, fmt.Errorf("invalid object id: %s", err.Error())
	}

	filter := bson.M{
		"type":      "Mental Health",
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endOfDay},
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

	daysInRange := int(endDate.Sub(startDate).Hours()/24) + 1

	dateIndexMap := make(map[string]int)
	for i := 0; i < daysInRange; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		dayIndex := int(startDate.AddDate(0, 0, i).Weekday())
		if _, exists := dateIndexMap[date]; !exists {
			dateIndexMap[date] = dayIndex
		}
	}

	dayData := make([]logEntry.MentalHealthRes, 7)

	for date, dayIndex := range dateIndexMap {
		idx := dayIndex % 7
		dayData[idx].Day = time.Weekday(idx).String()[:3]
		dayData[idx].Date = date
	}

	dayCount := make([]int, 7)
	for _, entry := range logEntryData {
		dayIndex := int(entry.When.Weekday())
		avg := dayData[dayIndex].AvgFeel

		dayData[dayIndex].AvgFeel = (avg*float64(dayCount[dayIndex]) + float64(mapFeel[entry.Feel])) / float64(dayCount[dayIndex]+1)
		dayData[dayIndex].Day = entry.When.Weekday().String()[:3]
		dayCount[dayIndex]++
	}

	responseData := make([]logEntry.MentalHealthRes, 0)
	dateSet := make(map[string]bool)
	for _, data := range dayData {
		if !dateSet[data.Date] {
			responseData = append(responseData, logEntry.MentalHealthRes{
				Date:    data.Date,
				Day:     data.Day,
				AvgFeel: math.Round(math.Abs(data.AvgFeel)*100) / 100,
			})
			dateSet[data.Date] = true
		}
	}

	sort.Slice(responseData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", responseData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", responseData[j].Date)
		return dateI.Before(dateJ)
	})

	// for i := range responseData {
	// 	if dayCount[i] != 0 {
	// 		responseData[i].AvgFeel *= 1
	// 	}
	// }

	return responseData, nil
}

func PhysicalHealthInsightWeeksData(c *fiber.Ctx) ([]logEntry.PhysicalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	weekStartDate := c.Query("firstDate")
	weekEndDate := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("Invalid startDate format: %s", err.Error())
	}
	endDate, err := time.Parse("2006-01-02", weekEndDate)
	if err != nil {
		return nil, fmt.Errorf("Invalid endDate format: %s", err.Error())
	}
	endOfDay := endDate.Add(24 * time.Hour)
	userId := c.Query("userId")
	userObjID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, fmt.Errorf("invalid object id: %s", err.Error())
	}

	filter := bson.M{
		"type":      "Physical Health",
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endOfDay},
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

	daysInRange := int(endDate.Sub(startDate).Hours()/24) + 1

	dateIndexMap := make(map[string]int)
	for i := 0; i < daysInRange; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		dayIndex := int(startDate.AddDate(0, 0, i).Weekday())
		if _, exists := dateIndexMap[date]; !exists {
			dateIndexMap[date] = dayIndex
		}
	}

	dayData := make([]logEntry.PhysicalHealthRes, 7)

	for date, dayIndex := range dateIndexMap {
		idx := dayIndex % 7
		dayData[idx].Day = time.Weekday(idx).String()[:3]
		dayData[idx].Date = date
	}

	dayCount := make([]int, 7)
	for _, entry := range logEntryData {
		dayIndex := int(entry.When.Weekday())
		avg := dayData[dayIndex].AvgPainLevel

		// avg := dayData[dayIndex].AvgPainLevel

		if dayCount[dayIndex] == 0 {
			dayData[dayIndex].AvgPainLevel = float64(entry.PainLevel)
		} else {
			sum := avg * float64(dayCount[dayIndex])
			dayData[dayIndex].AvgPainLevel = (sum + float64(entry.PainLevel)) / float64(dayCount[dayIndex]+1)
		}

		dayData[dayIndex].Day = entry.When.Weekday().String()[:3]
		dayCount[dayIndex]++
	}

	responseData := make([]logEntry.PhysicalHealthRes, 0)
	dateSet := make(map[string]bool)
	for _, data := range dayData {
		if !dateSet[data.Date] {
			responseData = append(responseData, logEntry.PhysicalHealthRes{
				Date:         data.Date,
				Day:          data.Day,
				AvgPainLevel: 0.0,
			})
			dateSet[data.Date] = true
		}
	}

	sort.Slice(responseData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", responseData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", responseData[j].Date)
		return dateI.Before(dateJ)
	})

	return responseData, nil
}
