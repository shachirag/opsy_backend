package userAuthenticate

import (
	"fmt"
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"opsy_backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Resend OTP
// @Description Resend 6 digit OTP to email
// @Tags user authorization
// @Accept application/json
// @Param user body userAuth.ResendOtpReqDto true "Resend 6 digit OTP to email"
// @Produce json
// @Success 200 {object} userAuth.UserPasswordResDto
// @Router /user/resend-otp [post]
func ResendOTP(c *fiber.Ctx) error {
	var (
		otpColl = database.GetCollection("otp")
		data    userAuth.ResendOtpReqDto
		otpData entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Find the existing OTP data for the provided email
	err = otpColl.FindOne(ctx, bson.M{"email": data.Email}).Decode(&otpData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(userAuth.UserPasswordResDto{
				Status:  false,
				Message: "No OTP found for the provided email. Please request a new OTP.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error while fetching OTP data: " + err.Error(),
		})
	}

	// Generate a new 6-digit OTP
	newOTP := utils.Generate6DigitOtp()

	// Update the existing OTP with the new OTP and reset the creation time
	otpData.Otp = newOTP
	otpData.CreatedAt = time.Now().UTC()

	update := bson.M{"$set": bson.M{"otp": newOTP, "createdAt": otpData.CreatedAt}}

	_, err = otpColl.UpdateOne(ctx, bson.M{"email": data.Email}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Failed to update OTP in the database: " + err.Error(),
		})
	}
	fmt.Println(newOTP)
	// Resending email to the recipient with the new OTP
	_, err = utils.SendEmail(data.Email, newOTP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error while resending email: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(userAuth.UserPasswordResDto{
		Status:  true,
		Message: "Successfully resent 6 digit OTP.",
	})
}
