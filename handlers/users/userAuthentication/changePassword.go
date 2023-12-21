package userAuthenticate

import (
	"opsy_backend/database"
	"opsy_backend/entity"

	userAuth "opsy_backend/dto/users/userAuthentication"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Change user Password
// @Description change admin Password
// @Tags user authorization
// @Accept application/json
// @Param id path string true "user ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param user body userAuth.ChangeUserPasswordReqDto true "Change password of user"
// @Produce json
// @Success 200 {object} userAuth.ChangeUserPasswordResDto
// @Router /user/change-password/{id} [put]
func ChangeUserPassword(c *fiber.Ctx) error {
	var (
		userColl = database.GetCollection("users")
		data     userAuth.ChangeUserPasswordReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	// Get the admin ID from the request parameters
	userID := c.Params("id")

	// Convert the admin ID string to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "Invalid user ID",
		})
	}

	// Find the admin document in MongoDB based on the provided admin ID
	filter := bson.M{"_id": objID}

	// Find the admin document in MongoDB
	result := userColl.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(userAuth.ChangeUserPasswordResDto{
				Status:  false,
				Message: "user not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "Error by finding user " + err.Error(),
		})
	}

	// Decode the admin data
	var user entity.UserEntity
	err = result.Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "Failed to decode user data: " + err.Error(),
		})
	}

	if data.CurrentPassword != "" {

		// Validate the current password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.CurrentPassword))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(userAuth.ChangeUserPasswordResDto{
				Status:  false,
				Message: "Current password is incorrect",
			})
		}
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "Failed to hash the new password",
		})
	}

	// Update the admin document with the new password
	update := bson.M{
		"$set": bson.M{
			"password": string(hashedNewPassword),
		},
	}

	// Execute the update operation
	updateRes, err := userColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "Failed to update user password in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(userAuth.ChangeUserPasswordResDto{
			Status:  false,
			Message: "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userAuth.ChangeUserPasswordResDto{
		Status:  true,
		Message: "user password updated successfully",
	})
}
