package userAuthenticate

import (
	"net/http"
	"opsy_backend/database"
	"opsy_backend/dto/users/logEntry"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary delete account
// @Tags user authorization
// @Description delete account
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "customer ID"
// @Produce json
// @Success 200 {object} logEntry.DeleteResDto
// @Router /user/delete-account/{id} [put]
func DeleteAccount(c *fiber.Ctx) error {
	customerID := c.Params("id")

	userObjectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		response := logEntry.DeleteResDto{
			Status:  false,
			Message: "Invalid user ID",
		}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	userColl := database.GetCollection("users")

	filter := bson.M{"_id": userObjectID}

	update := bson.M{"$set": bson.M{"isDeleted": true}}

	_, err = userColl.UpdateOne(ctx, filter, update)
	if err != nil {
		response := logEntry.DeleteResDto{
			Status:  false,
			Message: "Failed to delete account status: " + err.Error(),
		}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := logEntry.DeleteResDto{
		Status:  true,
		Message: "Account deleted successfully",
	}

	return c.Status(http.StatusOK).JSON(response)
}
