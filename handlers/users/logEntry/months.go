package logEntry

import (
	"fmt"
	"opsy_backend/database"
	logEntry "opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Months fetches log entry data for a specific month and year.
// @Summary Fetch log entry data for a specific month and year.
// @Tags logEntry
// @Description Fetch log entry data for a specific month and year.
// @Produce json
// @Param userId query string true "user ID"
// @Param month query string true "Month (01-12)"
// @Param year query string true "Year (YYYY)"
// @Success 200 {object} logEntry.CatgoriesResDto "Successfully fetched log entry data"
// @Failure 400 {object} logEntry.CatgoriesResDto "Invalid date format or missing parameters"
// @Failure 500 {object} logEntry.CatgoriesResDto "Failed to fetch or process data"
// @Router /user/months [get]
func Months(c *fiber.Ctx) error {
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
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	date, err := time.Parse("2006-01-02", yearStr+"-"+monthStr+"-01")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "Invalid date format: " + err.Error(),
		})
	}

	userId := c.Query("userId")
	userObjID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return c.Status(400).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}
	// Adjust the format of todayStart to match the createdAt format in MongoDB

	filter := bson.M{
		"isDeleted": false,
		"when":      bson.M{"$gte": date, "$lte": getLastDateOfMonth(date)},
		"userId":    userObjID,
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

func getLastDateOfMonth(t time.Time) time.Time {
	// Get the first day of the next month
	firstDayOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())

	// Subtract one day to get the last day of the current month
	lastDayOfMonth := firstDayOfNextMonth.Add(-time.Second)

	return lastDayOfMonth
}
