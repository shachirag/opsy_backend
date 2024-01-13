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
func FetchAllData(c *fiber.Ctx) error {
	// userId, err := auth.GetUserIdFromToken(c.Get("Authorization"))
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(logEntry.LogentryResDto{
	// 		Status:  false,
	// 		Message: "Unauthorized: " + err.Error(),
	// 	})
	// }
	var (
		logEntryColl = database.GetCollection("logEntry")
	)

	dateStr := c.Query("date")

	// Parse the string date to a time.Time type
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "Invalid date format: " + err.Error(),
		})
	}

	// Adjust the format of todayStart to match the createdAt format in MongoDB

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": date, "$lte": date.Add(24 * time.Hour)},
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	// Fetch data based on filters
	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		fmt.Println("Error fetching data:", err.Error()) // Log the error
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "Failed to fetch data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var logEntryData []logEntry.CategoriesRes
	for cursor.Next(ctx) {
		var logEntryEntity entity.LogEntryEntity
		if err := cursor.Decode(&logEntryEntity); err != nil {
			fmt.Println("Error decoding data:", err.Error()) // Log the decoding error
			return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CatgoriesResDto{
				Status:  false,
				Message: "Failed to decode data: " + err.Error(),
			})
		}
		logEntryData = append(logEntryData, logEntry.CategoriesRes{
			Id:          logEntryEntity.Id,
			Type:        logEntryEntity.Type,
			Notes:       logEntryEntity.Notes,
			Ways:        logEntryEntity.Ways,
			When:        logEntryEntity.When,
			PainLevel:   logEntryEntity.PainLevel,
			WhatItIsFor: logEntryEntity.WhatItIsFor,
			Feel:        logEntryEntity.Feel,
			Alert:       logEntryEntity.Alert,
			CreatedAt:   logEntryEntity.CreatedAt,
			UpdatedAt:   logEntryEntity.UpdatedAt,
		})
	}

	// Check if the result set is empty
	if len(logEntryData) == 0 {
		return c.Status(fiber.StatusOK).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "No data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.CatgoriesResDto{
		Status:  true,
		Message: "Successfully fetched the data.",
		Data:    logEntryData,
	})
}
