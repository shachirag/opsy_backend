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

var mapFeel = map[string]int{
	"In Crisis":  -1,
	"Struggling": -2,
	"Surviving":  -3,
	"Thriving":   -4,
	"Excelling":  -5,
}

func MentalHealthInsightWeeks(c *fiber.Ctx) error {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

	weekStartDate := c.Query("firstDate")
	weekEndDate := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", weekStartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.MentalHealthInsightResDto{
			Status:  false,
			Message: "Invalid startDate format: " + err.Error(),
		})
	}
	endDate, err := time.Parse("2006-01-02", weekEndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.MentalHealthInsightResDto{
			Status:  false,
			Message: "Invalid endDate format: " + err.Error(),
		})
	}
	endOfDay := endDate.Add(24 * time.Hour)

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": startDate, "$lte": endOfDay},
	}
	sortOptions := options.Find().SetSort(bson.M{"when": 1})

	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		fmt.Println("Error fetching data:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.MentalHealthInsightResDto{
			Status:  false,
			Message: "Failed to fetch data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var logEntryData []entity.LogEntryEntity
	if err := cursor.All(ctx, &logEntryData); err != nil {
		fmt.Println("Error decoding data:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.MentalHealthInsightResDto{
			Status:  false,
			Message: "Failed to decode data: " + err.Error(),
		})
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

	dayData := make([]logEntry.MentalHealthInsightRes, 7)

	for date, dayIndex := range dateIndexMap {
		idx := dayIndex % 7
		dayData[idx].Day = time.Weekday(idx).String()[:3]
		dayData[idx].Date = date
	}

	dayCount := make([]int, 7)
	for _, entry := range logEntryData {
		dayIndex := int(entry.When.Weekday())
		avg := dayData[dayIndex].AvgFeel

		if avg == 0 {
			dayData[dayIndex].AvgFeel = float64(mapFeel[entry.Feel])
		} else {
			dayData[dayIndex].AvgFeel = (avg*float64(dayCount[dayIndex]) + float64(mapFeel[entry.Feel])) / float64(dayCount[dayIndex]+1)
		}

		dayData[dayIndex].Day = entry.When.Weekday().String()[:3]
		dayCount[dayIndex]++
	}

	responseData := make([]logEntry.MentalHealthInsightRes, 0)
	dateSet := make(map[string]bool)
	for _, data := range dayData {
		if !dateSet[data.Date] {
			responseData = append(responseData, data)
			dateSet[data.Date] = true
		}
	}

	sort.Slice(responseData, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", responseData[i].Date)
		dateJ, _ := time.Parse("2006-01-02", responseData[j].Date)
		return dateI.Before(dateJ)
	})

	// Adjust the precision of float values
	for i := range responseData {
		responseData[i].AvgFeel = math.Round(math.Abs(responseData[i].AvgFeel)*100) / 100
		if dayCount[i] != 0 {
			responseData[i].AvgFeel *= -1
		}
	}

	response := logEntry.MentalHealthInsightResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data:    responseData,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
