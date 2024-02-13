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

// @Summary Update mentalHealth
// @Description Update mentalHealth
// @Tags logEntry
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "mentalHealth ID"
// @Param logEntry formData logEntry.EditMentalHealthReqDto true "Update data of mentalHealth"
// @Produce json
// @Success 200 {object} logEntry.EditMentalHealthResDto
// @Router /user/edit-mentalHealth/{id} [put]
func UpdateMentalHealth(c *fiber.Ctx) error {
	ctx := c.Context()

	var (
		logentryColl = database.GetCollection("logEntry")
		logentry     entity.LogEntryEntity
		data         logEntry.EditMentalHealthReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if the appointment ID is provided in the request
	mentalHealthID := c.Params("id")
	if mentalHealthID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: "mentalHealth ID is missing in the request",
		})
	}

	objID, err := primitive.ObjectIDFromHex(mentalHealthID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: "Invalid mentalHealth ID",
		})
	}

	filter := bson.M{"_id": objID, "type": "Mental Health"}
	result := logentryColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(logEntry.EditMentalHealthResDto{
				Status:  false,
				Message: "mentalHealth not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: "Internal server error: " + err.Error(),
		})
	}

	// Decode the appointment data
	err = result.Decode(&logentry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: "Failed to decode mentalHealth data: " + err.Error(),
		})
	}

	dateTime, err := time.Parse("2006-01-02T15:04:05", data.When)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Update the appointment document with new data
	update := bson.M{
		"$set": bson.M{
			"when":      dateTime,
			"notes":     data.Notes,
			"feel":      data.Feel,
			"ways":      data.Ways,
			"updatedAt": time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := logentryColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: "Failed to update mentalHealth data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(logEntry.EditMentalHealthResDto{
			Status:  false,
			Message: "mentalHealth not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.EditMentalHealthResDto{
		Status:  true,
		Message: "mentalHealth data updated successfully",
	})
}
