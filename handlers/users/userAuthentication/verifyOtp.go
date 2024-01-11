package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Verify user OTP
// @Description Verify the entered 6 digit OTP
// @Tags user authorization
// @Accept application/json
// @Param user body userAuth.VerifyOtpReqDto true "Verify 6 digit OTP"
// @Produce json
// @Success 200 {object} userAuth.UserPasswordResDto
// @Router /user/verify-otp [post]
func VerifyOtp(c *fiber.Ctx) error {
	var (
		otpColl = database.GetCollection("otp")
		data    userAuth.VerifyOtpReqDto
		otpData entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Entered OTP is required",
		})
	}
	smallEmail := strings.ToLower(data.Email)
	// Find the user with email address from client
	err = otpColl.FindOne(ctx, bson.M{"email": smallEmail}, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&otpData)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(userAuth.UserPasswordResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}

		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the user: " + err.Error(),
		})
	}

	// Compare the entered OTP with the OTP from the database
	if data.EnteredOTP != otpData.Otp {

		return c.Status(400).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	return c.Status(200).JSON(userAuth.UserPasswordResDto{
		Status:  true,
		Message: "OTP verified successfully",
	})
}
