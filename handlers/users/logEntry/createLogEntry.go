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
	"go.mongodb.org/mongo-driver/mongo/options"
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
		logEntryColl       = database.GetCollection("logEntry")
		yearlyInsightsColl = database.GetCollection("yearlyInsights")
		data               logEntry.LogEntryReqDto
	)
	// Parsing the request body
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

	logEntryData := entity.LogEntryEntity{
		Id:          id,
		Type:        data.Type,
		IsDeleted:   false,
		Feel:        data.Feel,
		Notes:       data.Notes,
		Ways:        data.Ways,
		When:        dateTime,
		PainLevel:   data.PainLevel,
		WhatItIsFor: data.WhatItIsFor,
		Alert:       data.Alert,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_, err = logEntryColl.InsertOne(ctx, logEntryData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
			Status:  false,
			Message: "failed to store Log Entry Database " + err.Error(),
		})
	}

	year := dateTime.Year()
	month := int32(dateTime.Month())

	filter := bson.M{"year": year, "month": month}
	update := bson.M{
		"$push": bson.M{
			"mentalHealth": bson.M{
				"avgFeel":              logEntryData.Feel,
				"totalMentalHealthLog": 1,
			},
			"physicalHealth": bson.M{
				"avgFeel":              logEntryData.Feel,
				"totalMentalHealthLog": 1,
			},
		},
	}

	_, err = yearlyInsightsColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {

		return c.Status(fiber.StatusOK).JSON(logEntry.GetLogEntryResDto{
			Status:  true,
			Message: "Data inserted successfully",
			Id:      id,
		})
	}
	return c.Status(fiber.StatusOK).JSON(logEntry.GetLogEntryResDto{
		Status:  true,
		Message: "Data inserted successfully",
		Id:      id,
	})
}
