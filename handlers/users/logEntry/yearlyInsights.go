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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch yearly insights data
// @Description Retrieves mental health and physical health data for a specified year
// @Tags logEntry
// @Accept json
// @Produce json
// @Param year query int true "Year (YYYY)"
// @Success 200 {object} YearlyInsightsResDto
// @Router /user/yearly-insights [get]

func YearlyInsights(c *fiber.Ctx) error {
	year := c.Query("year")

	// Validate the year format
	if _, err := time.Parse("2006", year); err != nil {
		return fmt.Errorf("Invalid year format: %s", err.Error())
	}

	startDate := time.Date(int(year), 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(1, 0, 0).Add(-time.Nanosecond)

	mentalHealthData, err := MentalHealthInsightYearData(c, startDate, endDate)
	if err != nil {
		return err
	}

	physicalHealthData, err := PhysicalHealthInsightYearData(c, startDate, endDate)
	if err != nil {
		return err
	}

	allMonths := make(map[string]bool)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
		allMonths[d.Format("01-2006")] = true
	}

	// Fill in missing data with 0 values for mental health
	for _, data := range mentalHealthData {
		delete(allMonths, data.TotalMentalHealthLog)
	}
	for month := range allMonths {
		mentalHealthData = append(mentalHealthData, logEntry.YearlyMentalHealthRes{
			TotalMentalHealthLog: month,
			AvgFeel:              0,
		})
	}

	// Fill in missing data with 0 values for physical health
	for _, data := range physicalHealthData {
		delete(allMonths, data.TotalPhysicalHealthLog)
	}
	for month := range allMonths {
		physicalHealthData = append(physicalHealthData, logEntry.YearlyPhysicalHealthRes{
			TotalPhysicalHealthLog: month,
			AvgFeel:                0,
		})
	}

	sort.Slice(mentalHealthData, func(i, j int) bool {
		dateI, _ := time.Parse("01-2006", mentalHealthData[i].TotalMentalHealthLog)
		dateJ, _ := time.Parse("01-2006", mentalHealthData[j].TotalMentalHealthLog)
		return dateI.Before(dateJ)
	})

	sort.Slice(physicalHealthData, func(i, j int) bool {
		dateI, _ := time.Parse("01-2006", physicalHealthData[i].TotalPhysicalHealthLog)
		dateJ, _ := time.Parse("01-2006", physicalHealthData[j].TotalPhysicalHealthLog)
		return dateI.Before(dateJ)
	})

	response := logEntry.YearlyInsightsResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: YearlyInsightsData{
			MentalHealth:   mentalHealthData,
			PhysicalHealth: physicalHealthData,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func MentalHealthInsightYearData(c *fiber.Ctx, startDate, endDate time.Time) ([]logEntry.YearlyMentalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

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

	var logEntryData []entity.YearlyInsightsEntity
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

	var yearlyData []logEntry.YearlyMentalHealthRes

	for date, values := range dailyData {
		dateTime, _ := time.Parse("2006-01-02", date)
		avg := 0.0
		totalDays := float64(dayCount[date])
		if totalDays > 0 {
			for _, value := range values {
				avg += value
			}
			avg /= totalDays
		}

		yearlyData = append(yearlyData, logEntry.YearlyMentalHealthRes{
			TotalMentalHealthLog: date,
			AvgFeel:              math.Round(math.Abs(avg)*100) / 100,
		})
	}

	return yearlyData, nil
}

func PhysicalHealthInsightYearData(c *fiber.Ctx, startDate, endDate time.Time) ([]logEntry.YearlyMentalHealthRes, error) {
	var (
		logEntryColl = database.GetCollection("logEntry")
		ctx          = c.Context()
	)

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

	var logEntryData []entity.YearlyInsightsEntity
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

	var yearlyData []logEntry.YearlyMentalHealthRes

	// Calculate average pain level for the entire year
	yearlyAvg := 0.0
	for _, values := range dailyData {
		for _, value := range values {
			yearlyAvg += value
		}
	}
	if totalDays := float64(len(dailyData) * 7); totalDays > 0 {
		yearlyAvg /= totalDays
	}

	// Append the result to the yearlyData slice
	yearlyData = append(yearlyData, logEntry.YearlyMentalHealthRes{
		TotalMentalHealthLog: startDate.Year(),
		AvgFeel:              yearlyAvg,
	})

	return yearlyData, nil
}
