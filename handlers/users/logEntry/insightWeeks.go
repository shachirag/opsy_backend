package logEntry

import (
	"fmt"
	"opsy_backend/database"
	logEntry "opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary fetch all required data
// @Tags logEntry
// @Description fetch all required data
// @Produce json
// @Param date query string true "date (YYYY-MM-DD)"
// @Success 200 {object} logEntry.CatgoriesResDto
// @Router /user/fetch-all-data [get]
func InsightWeeks(c *fiber.Ctx) error {
	var (
		logEntryColl = database.GetCollection("logEntry")
	)
	weekStartDate := c.Query("firstDate")
	weekEndDate := c.Query("endDate")
	startDate, err := time.Parse("2006-01-02", weekStartDate)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(logEntry.InsightResDto{
            Status:  false,
            Message: "Invalid startDate format: " + err.Error(),
        })
    }
    endDate, err := time.Parse("2006-01-02", weekEndDate)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(logEntry.InsightResDto{
            Status:  false,
            Message: "Invalid endDate format: " + err.Error(),
        })
    }
	endOfDay := endDate.Add(24 * time.Hour)
	// Adjust the format of todayStart to match the createdAt format in MongoDB

	filter := bson.M{
		"isDeleted": false,
		"when": bson.M{"$gte": startDate, "$lte": endOfDay },
	}
	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})
	// Fetch data based on filters
	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		fmt.Println("Error fetching data:", err.Error()) // Log the error
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.InsightResDto{
			Status:  false,
			Message: "Failed to fetch data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var logEntryData []logEntry.InsightRes
	for cursor.Next(ctx) {
		var logEntryEntity entity.LogEntryEntity
		if err := cursor.Decode(&logEntryEntity); err != nil {
			fmt.Println("Error decoding data:", err.Error()) // Log the decoding error
			return c.Status(fiber.StatusInternalServerError).JSON(logEntry.InsightResDto{
				Status:  false,
				Message: "Failed to decode data: " + err.Error(),
			})
		}
		logEntryData = append(logEntryData, logEntry.InsightRes{
			Id:        logEntryEntity.Id,
			Type:      logEntryEntity.Type,
			Ways:      logEntryEntity.Ways,
			When:      logEntryEntity.When,
			PainLevel: logEntryEntity.PainLevel,
			CreatedAt: logEntryEntity.CreatedAt,
			UpdatedAt: logEntryEntity.UpdatedAt,
		})
	}

	// Check if the result set is empty
	if len(logEntryData) == 0 {
		return c.Status(fiber.StatusOK).JSON(logEntry.InsightResDto{
			Status:  false,
			Message: "No data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.InsightResDto{
		Status:  true,
		Message: "Successfully fetched the data.",
		Data:    logEntryData,
	})
}

