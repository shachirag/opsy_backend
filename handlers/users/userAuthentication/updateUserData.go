package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update user
// @Description Update user
// @Tags user authorization
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "user ID"
// @Param user formData userAuth.UpdateUserReqDto true "Update data of user"
// @Produce json
// @Success 200 {object} userAuth.UpdateUserResDto
// @Router /user/update-user-data/{id} [put]
func UpdateUser(c *fiber.Ctx) error {

	var (
		userColl = database.GetCollection("users")
		customer entity.UserEntity
		data     userAuth.UpdateUserReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if the user ID is provided in the request
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: "user ID is missing in the request",
		})
	}

	// Find the user document in MongoDB based on the provided user ID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: "Invalid user ID",
		})
	}

	// Find the user document in MongoDB
	filter := bson.M{"_id": objID}
	result := userColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(userAuth.UpdateUserResDto{
				Status:  false,
				Message: "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: "internal server error " + err.Error(),
		})
	}

	// Update the admin document with new data
	update := bson.M{
		"name":      data.Name,
		"updatedAt": time.Now().UTC(),
	}

	updateFields := bson.M{"$set": update}
	// Execute the update operation
	updateRes, err := userColl.UpdateOne(ctx, filter, updateFields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: "Failed to update user data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: "user not found",
		})
	}

	err = result.Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UpdateUserResDto{
			Status:  false,
			Message: "Failed to decode updated user data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(userAuth.UpdateUserResDto{
		Status:  true,
		Message: "user data updated successfully",
	})
}
