package logEntry

import (
	"opsy_backend/database"
	"opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch logentry By ID
// @Description Fetch logentry By ID
// @Tags logEntry
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "logentry ID"
// @Produce json
// @Success 200 {object} logEntry.LogentryResDto
// @Router /user/logentry-info/{id} [get]
func FetchLogEntryById(c *fiber.Ctx) error {
	// userId, err := auth.GetUserIdFromToken(c.Get("Authorization"))
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(logEntry.LogentryResDto{
	// 		Status:  false,
	// 		Message: "Unauthorized: " + err.Error(),
	// 	})
	// }
	var logentry entity.LogEntryEntity

	logentryId := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(logentryId)

	if err != nil {
		return c.Status(400).JSON(logEntry.LogentryResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})

	}

	logentryColl := database.GetCollection("logEntry")

	err = logentryColl.FindOne(ctx, bson.M{"_id": objId,}).Decode(&logentry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(logEntry.LogentryResDto{
				Status:  false,
				Message: "logentry not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.LogentryResDto{
			Status:  false,
			Message: "Failed to fetch logentry from MongoDB: " + err.Error(),
		})
	}

	logentryRes := logEntry.LogentryRes{
		Id:          logentry.Id,
		Type:        logentry.Type,
		Feel:        logentry.Feel,
		Notes:       logentry.Notes,
		Ways:        logentry.Ways,
		When:        logentry.When,
		PainLevel:   logentry.PainLevel,
		WhatItIsFor: logentry.WhatItIsFor,
		Alert:       logentry.Alert,
		CreatedAt:   logentry.CreatedAt,
		UpdatedAt:   logentry.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.LogentryResDto{
		Status:  true,
		Message: "logentry data retrieved successfully",
		Data:    &logentryRes,
	})
}
