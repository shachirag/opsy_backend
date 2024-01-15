package logEntry

import (
	"context"
	"opsy_backend/database"
	logEntry "opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()

// @Summary CreateLogEntry
// @Description  CreateLogEntry
// @Tags logEntry
// @Param user body logEntry.LogEntryReqDto true "CreateLogEntry for user"
// @Produce json
// @Success 200 {object} logEntry.GetLogEntryResDto
// @Router /user/create-log-entry [post]
func CreateLogEntry(c *fiber.Ctx) error {
	var (
		logEntryColl = database.GetCollection("logEntry")
		data         logEntry.LogEntryReqDto
		logEntryData entity.LogEntryEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(logEntry.GetLogEntryResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	dateTime, err := time.Parse("2006-01-02T15:04:05", data.DateTime)
	if err != nil {
		return c.Status(500).JSON(logEntry.GetLogEntryResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	id := primitive.NewObjectID()

	count, err := logEntryColl.CountDocuments(ctx, bson.M{"isDeleted": false})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
				Status:  false,
				Message: "Error getting count : " + err.Error(),
			})
		}
	}

	logEntryData = entity.LogEntryEntity{
		Id:          id,
		Type:        data.Type,
		IsDeleted:   false,
		Feel:        data.Feel,
		Notes:       data.Notes,
		Ways:        data.Ways,
		When:        dateTime,
		PainLevel:   data.PainLevel,
		WhatItIsFor: data.WhatItIsFor,
		UserId:      data.UserId,
		Alert:       data.Alert,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		NumberCount: count + 1,
	}

	_, err = logEntryColl.InsertOne(ctx, logEntryData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
			Status:  false,
			Message: "failed to store Log Entry Database " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(logEntry.GetLogEntryResDto{
		Status:  true,
		Message: "Data inserted successfully",
		Data: logEntry.LogEntryRes{
			Id:          id,
			Type:        data.Type,
			Feel:        data.Feel,
			Notes:       data.Notes,
			Ways:        data.Ways,
			When:        dateTime,
			PainLevel:   data.PainLevel,
			WhatItIsFor: data.WhatItIsFor,
			Alert:       data.Alert,
			NumberCount: count + 1,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	})
}
