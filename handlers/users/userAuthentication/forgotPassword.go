package userAuthenticate

import (
	"fmt"
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"opsy_backend/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Forgot Password
// @Description Forgot Password
// @Tags user authorization
// @Accept application/json
// @Param user body userAuth.UserForgotPasswordReqDto true "forgot password for user"
// @Produce json
// @Success 200 {object} userAuth.UserPasswordResDto
// @Router /user/forgot-password [post]
func ForgotPassword(c *fiber.Ctx) error {
	var (
		userColl = database.GetCollection("users")
		otpColl  = database.GetCollection("otp")
		data     userAuth.UserForgotPasswordReqDto
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
				Message: "We couldn’t find an account with that email address.",
			})
		}

		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the user: " + err.Error(),
		})
	}

	// Generate 6-digit OTP
	otp := utils.Generate6DigitOtp()
	fmt.Println(otp)
	// Store the OTP in the forgotPassword collection
	otpData := entity.OtpEntity{
		Id:        primitive.NewObjectID(),
		Otp:       otp,
		Email:     smallEmail,
		CreatedAt: time.Now().UTC(),
	}

	_, err = otpColl.InsertOne(ctx, otpData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Failed to store OTP in the database: " + err.Error(),
		})
	}

	// Sending email to the recipient with the OTP
	_, err = utils.SendEmail(data.Email, user.Name, otp)
	if err != nil {
		return c.Status(500).JSON(userAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error, while sending email: " + err.Error(),
		})
	}

	return c.Status(200).JSON(userAuth.UserPasswordResDto{
		Status:  true,
		Message: "Successfully send 6 digit OTP.",
	})
}
