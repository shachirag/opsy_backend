package logEntry

import (
	"context"
	"opsy_backend/database"
	logEntry "opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		// yearlyInsightsColl = database.GetCollection("yearlyInsight")
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

	// var mapFeel = map[string]int{
	// 	"In Crisis":  -1,
	// 	"Struggling": -2,
	// 	"Surviving":  -3,
	// 	"Thriving":   -4,
	// 	"Excelling":  -5,
	// }

	// // Calculate the weighted average feel based on mapFeel values
	// weightedAvgFeel := mapFeel[data.Feel]

	// // Update Yearly Insights
	// year, month, _ := logEntryData.When.Date()
	// filter := bson.M{"year": year, "month": int32(month)}

	// mentalHealthUpdate := bson.M{
	// 	"$inc": bson.M{
	// 		"mentalHealth.$.totalMentalHealthLog": 1,
	// 		"mentalHealth.$.avgFeel":              weightedAvgFeel,
	// 	},
	// }

	// physicalHealthUpdate := bson.M{
	// 	"$inc": bson.M{
	// 		"physicalHealth.$.totalMentalHealthLog": 1,
	// 		"physicalHealth.$.avgPain":              logEntryData.PainLevel,
	// 	},
	// }

	// result, err := yearlyInsightsColl.UpdateOne(ctx, filter, mentalHealthUpdate)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
	// 		Status:  false,
	// 		Message: "failed to update Yearly Insights " + err.Error(),
	// 	})
	// }

	// // Update physical health
	// _, err = yearlyInsightsColl.UpdateOne(ctx, filter, physicalHealthUpdate)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
	// 		Status:  false,
	// 		Message: "failed to update Yearly Insights (physical health) " + err.Error(),
	// 	})
	// }

	// if result.MatchedCount == 0 {
	// 	// If the document does not exist, create a new one
	// 	yearlyInsightsData := entity.YearlyInsightsEntity{
	// 		Month: int32(month),
	// 		Year:  int32(year),
	// 		MentalHealth: []entity.MentalHealth{
	// 			{
	// 				AvgFeel:              weightedAvgFeel,
	// 				TotalMentalHealthLog: 1,
	// 			},
	// 		},
	// 		PhysicalHealth: []entity.PhysicalHealth{
	// 			{
	// 				AvgPain:              logEntryData.PainLevel,
	// 				TotalMentalHealthLog: 1,
	// 			},
	// 		},
	// 	}

	// 	_, err := yearlyInsightsColl.InsertOne(ctx, yearlyInsightsData)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
	// 			Status:  false,
	// 			Message: "failed to insert new Yearly Insights document " + err.Error(),
	// 		})
	// 	}
	// }

	return c.Status(fiber.StatusOK).JSON(logEntry.GetLogEntryResDto{
		Status:  true,
		Message: "Data inserted successfully",
		Id:      id,
	})
}
