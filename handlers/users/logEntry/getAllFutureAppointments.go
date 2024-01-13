package logEntry

import (
	"encoding/json"
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

// @Summary fetch all required data
// @Tags logEntry
// @Description fetch all required data
// @Produce json
// @Success 200 {object} logEntry.FutureAppointmentDto
// @Router /user/get-future-appointments [get]
func FetchFutureAppointments(c *fiber.Ctx) error {

	var (
		logEntryColl = database.GetCollection("logEntry")
	)

	userId := c.Query("userId")
	userObjID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return c.Status(400).JSON(logEntry.FutureAppointmentDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	// Format the current time using "2006-01-02T15:04:05" layout
	currentTimeFormatted := time.Now().UTC().Format("2006-01-02T15:04:05")

	filter := bson.M{
		"isDeleted": false,
		"userId":    userObjID,
		"when":      bson.M{"$gt": currentTimeFormatted},
	}

	fmt.Println(filter)

	data, _ := json.Marshal(filter)
	fmt.Println(string(data))

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	// Fetch data based on filters
	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		fmt.Println("Error fetching data:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.FutureAppointmentDto{
			Status:  false,
			Message: "Failed to fetch data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var logEntryData []logEntry.FutureAppointmentRes
	for cursor.Next(ctx) {
		var logEntryEntity entity.LogEntryEntity
		if err := cursor.Decode(&logEntryEntity); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(logEntry.FutureAppointmentDto{
				Status:  false,
				Message: "Failed to decode data: " + err.Error(),
			})
		}
		logEntryData = append(logEntryData, logEntry.FutureAppointmentRes{
			Id:          logEntryEntity.Id,
			Notes:       logEntryEntity.Notes,
			When:        logEntryEntity.When,
			WhatItIsFor: logEntryEntity.WhatItIsFor,
			Alert:       logEntryEntity.Alert,
			UpdatedAt:   logEntryEntity.UpdatedAt,
			CreatedAt:   logEntryEntity.CreatedAt,
		})
	}

	// Check if the result set is empty
	if len(logEntryData) == 0 {
		return c.Status(fiber.StatusOK).JSON(logEntry.FutureAppointmentDto{
			Status:  false,
			Message: "No data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.FutureAppointmentDto{
		Status:  true,
		Message: "Successfully fetched the data.",
		Data:    logEntryData,
	})
}
