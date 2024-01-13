package logEntry

import (
	"fmt"
	"math"
	"opsy_backend/database"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var mapFeel = map[string]int{
// 	"In Crisis":  -1,
// 	"Struggling": -2,
// 	"Surviving":  -3,
// 	"Thriving":   -4,
// 	"Excelling":  -5,
// }

// @Summary Fetch monthly insights data for a specific year
// @Description Retrieves mental health and physical health data for each month of the specified year
// @Tags logEntry
// @Accept json
// @Produce json
// @Param year query string true "Year (YYYY)"
// @Success 200 {object} YearInsightsResDto
// @Router /user/yearly-insights [get]
func YearlyInsights(c *fiber.Ctx) error {

	year := c.Query("year")
	parsedYear, err := time.Parse("2006", year)
	if err != nil {
		return fmt.Errorf("Invalid year format: %s", err.Error())
	}

	var (
		months             = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
		mentalHealthData   = make([]YearlyMentalHealthInsightsResData, 0, 12)
		physicalHealthData = make([]YearlyPhysicalHealthInsightsResData, 0, 12)
	)

	for i := 1; i <= 12; i++ {
		monthStartDate := time.Date(parsedYear.Year(), time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEndDate := monthStartDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

		mentalData, err := MentalHealthInsightMonthsData1(c, monthStartDate, monthEndDate)
		if err != nil {
			return err
		}
		physicalData, err := PhysicalHealthInsightMonthsData1(c, monthStartDate, monthEndDate)
		if err != nil {
			return err
		}

		monthlyMentalData := createEmptyMonthlyMentalData(months[i-1])
		monthlyPhysicalData := createEmptyMonthlyPhysicalData(months[i-1])

		if len(mentalData) > 0 {
			monthlyMentalData.AvgFeel = math.Round(math.Abs(mentalData[0].AvgFeel)*100) / 100
		}

		if len(physicalData) > 0 {
			monthlyPhysicalData.AvgPainLevel = physicalData[0].AvgPainLevel
		}

		mentalHealthData = append(mentalHealthData, monthlyMentalData)
		physicalHealthData = append(physicalHealthData, monthlyPhysicalData)
	}

	response := YearInsightsResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: YearlyInsightsRes{
			MentalHealth:   mentalHealthData,
			PhysicalHealth: physicalHealthData,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func createEmptyMonthlyMentalData(month string) YearlyMentalHealthInsightsResData {
	return YearlyMentalHealthInsightsResData{
		Month:   month,
		AvgFeel: 0,
	}
}

func createEmptyMonthlyPhysicalData(month string) YearlyPhysicalHealthInsightsResData {
	return YearlyPhysicalHealthInsightsResData{
		Month:        month,
		AvgPainLevel: 0,
	}
}

func PhysicalHealthInsightMonthsData1(c *fiber.Ctx, startDate, endDate time.Time) ([]YearlyPhysicalHealthInsightsResData, error) {
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

	var logEntryData []entity.LogEntryEntity
	if err := cursor.All(ctx, &logEntryData); err != nil {
		return nil, fmt.Errorf("Failed to decode data: %s", err.Error())
	}

	dailyData := make(map[string][]float64)
	dayCount := make(map[string]int)

	for _, entry := range logEntryData {
		dateKey := entry.When.Format("2006-01-02")
		dayIndex := int(entry.When.Weekday())

		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = make([]float64, 7)
		}

		dailyData[dateKey][dayIndex] += float64(entry.PainLevel)
		dayCount[dateKey]++
	}

	var monthlyData []YearlyPhysicalHealthInsightsResData

	for date, values := range dailyData {
		dateTime, _ := time.Parse("2006-01-02", date)
		avg := 0.0

		for _, value := range values {
			avg += value
		}

		avg /= float64(dayCount[date])
		monthlyData = append(monthlyData, YearlyPhysicalHealthInsightsResData{
			Month:        dateTime.Month().String(),
			AvgPainLevel: avg,
		})
	}

	return monthlyData, nil
}

func MentalHealthInsightMonthsData1(c *fiber.Ctx, startDate, endDate time.Time) ([]YearlyMentalHealthInsightsResData, error) {
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

	var logEntryData []entity.LogEntryEntity
	if err := cursor.All(ctx, &logEntryData); err != nil {
		return nil, fmt.Errorf("Failed to decode data: %s", err.Error())
	}

	dailyData := make(map[string][]float64)
	dayCount := make(map[string]int)

	for _, entry := range logEntryData {
		dateKey := entry.When.Format("2006-01-02")
		dayIndex := int(entry.When.Weekday())

		if _, ok := dailyData[dateKey]; !ok {
			dailyData[dateKey] = make([]float64, 7)
		}

		dailyData[dateKey][dayIndex] += float64(mapFeel[entry.Feel])
		dayCount[dateKey]++
	}

	var monthlyData []YearlyMentalHealthInsightsResData

	for date, values := range dailyData {
		dateTime, _ := time.Parse("2006-01-02", date)
		avg := 0.0

		for _, value := range values {
			avg += value
		}

		avg /= float64(dayCount[date])
		monthlyData = append(monthlyData, YearlyMentalHealthInsightsResData{
			Month:   dateTime.Month().String(),
			AvgFeel: math.Round(math.Abs(avg)*100) / 100,
		})
	}

	return monthlyData, nil
}

type YearlyMentalHealthInsightsResData struct {
	Month   string  `json:"month"`
	AvgFeel float64 `json:"avgFeel"`
}

type YearlyPhysicalHealthInsightsResData struct {
	Month        string  `json:"month"`
	AvgPainLevel float64 `json:"avgPainLevel"`
}

type YearlyInsightsRes struct {
	MentalHealth   []YearlyMentalHealthInsightsResData   `json:"mentalHealth"`
	PhysicalHealth []YearlyPhysicalHealthInsightsResData `json:"physicalHealth"`
}

type YearInsightsResDto struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    YearlyInsightsRes `json:"data"`
}
