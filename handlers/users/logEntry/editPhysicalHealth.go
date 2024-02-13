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

// @Summary Update physicalHealth
// @Description Update physicalHealth
// @Tags logEntry
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "physicalHealth ID"
// @Param logEntry formData logEntry.EditPhysicalHealthReqDto true "Update data of physicalHealth"
// @Produce json
// @Success 200 {object} logEntry.EditPhysicalHealthResDto
// @Router /user/edit-physicalHealth/{id} [put]
func UpdatePhysicalHealth(c *fiber.Ctx) error {
	ctx := c.Context()

	var (
		logentryColl = database.GetCollection("logEntry")
		logentry     entity.LogEntryEntity
		data         logEntry.EditPhysicalHealthReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if the appointment ID is provided in the request
	physicalHealthID := c.Params("id")
	if physicalHealthID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: "physicalHealth ID is missing in the request",
		})
	}

	objID, err := primitive.ObjectIDFromHex(physicalHealthID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: "Invalid physicalHealth ID",
		})
	}

	filter := bson.M{"_id": objID,"type":"Physical Health"}
	result := logentryColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(logEntry.EditPhysicalHealthResDto{
				Status:  false,
				Message: "physicalHealth not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: "Internal server error: " + err.Error(),
		})
	}

	// Decode the physicalHealth data
	err = result.Decode(&logentry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: "Failed to decode physicalHealth data: " + err.Error(),
		})
	}

	dateTime, err := time.Parse("2006-01-02T15:04:05", data.When)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Update the physicalHealth document with new data
	update := bson.M{
		"$set": bson.M{
			"when":      dateTime,
			"notes":     data.Notes,
			"painlevel": data.PainLevel,
			"ways":      data.Ways,
			"updatedAt": time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := logentryColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: "Failed to update physicalHealth data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(logEntry.EditPhysicalHealthResDto{
			Status:  false,
			Message: "physicalHealth not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.EditPhysicalHealthResDto{
		Status:  true,
		Message: "physicalHealth data updated successfully",
	
	})
}
