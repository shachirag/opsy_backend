package logEntry

import (
	"opsy_backend/database"
	"opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update appointment
// @Description Update appointment
// @Tags logEntry
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "appointment ID"
// @Param logEntry formData logEntry.AppointmentReqDto true "Update data of appointment"
// @Produce json
// @Success 200 {object} logEntry.GetAppointmentResDto
// @Router /user/edit-appointment/{id} [put]
func UpdateLogEntry(c *fiber.Ctx) error {

	var (
		logentryColl = database.GetCollection("logEntry")
		logentry     entity.LogEntryEntity
		data         logEntry.AppointmentReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if the customer ID is provided in the request
	appointmentID := c.Params("id")
	if appointmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: "appointment ID is missing in the request",
		})
	}

	objID, err := primitive.ObjectIDFromHex(appointmentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{"_id": objID}
	result := logentryColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(logEntry.GetAppointmentResDto{
				Status:  false,
				Message: "appointment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: "internal server error " + err.Error(),
		})
	}

	// Decode the customer data
	err = result.Decode(&logentry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: "Failed to decode appointment data: " + err.Error(),
		})
	}

	dateTime, err := time.Parse("2006-01-02T15:04:05", data.When)
	if err != nil {
		return c.Status(500).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Update the admin document with new data
	update := bson.M{
		"$set": bson.M{
			"when":        dateTime,
			"whatItIsFor": data.WhatItIsFor,
			"alert":       data.Alert,
			"notes":       data.Notes,
			"updatedAt":   time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := logentryColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: "Failed to update appointment data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(logEntry.GetAppointmentResDto{
			Status:  false,
			Message: "appointment not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.GetAppointmentResDto{
		Status:  true,
		Message: "appointment data updated successfully",
	})
}
