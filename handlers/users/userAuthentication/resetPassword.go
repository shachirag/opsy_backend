package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"opsy_backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Reset user Password after OTP Verification
// @Description Reset user password after OTP verification using the new password and confirm password
// @Tags user authorization
// @Accept application/json
// @Param user body userAuth.ResetPasswordAfterOtpReqDto true "Reset user password after OTP verification"
// @Produce json
// @Success 200 {object} userAuth.UserPasswordResDto
// @Router /user/reset-password [put]
func ResetPasswordAfterOtp(c *fiber.Ctx) error {
	var (
		userColl = database.GetCollection("users")
		data     userAuth.ResetPasswordAfterOtpReqDto
		user     entity.UserEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}
	smallEmail := strings.ToLower(data.Email)
	// Find the user with email address from client
	err = userColl.FindOne(ctx, bson.M{"email": smallEmail, "isDeleted": false}).Decode(&user)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(userAuth.UserPasswordResDto{
				Status:  false,
				Message: "No user found",
			})
		}

		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the user: " + err.Error(),
		})
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Invalid Crdentials " + err.Error(),
		})
	}

	// Update the user's password in the database
	_, err = userColl.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Failed to update password in the database: " + err.Error(),
		})
	}

	return c.Status(200).JSON(userAuth.UserPasswordResDto{
		Status:  true,
		Message: "Password updated successfully after OTP verification",
	})
}
