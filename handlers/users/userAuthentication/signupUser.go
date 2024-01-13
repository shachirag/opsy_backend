package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"opsy_backend/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Signup user
// @Description Signup user
// @Tags user authorization
// @Param signup body userAuth.SignupReqDto true "send 6 digit otp to email for signup"
// @Produce json
// @Success 200 {object} userAuth.SignupResDto
// @Router /user/signup [post]
func SignupUser(c *fiber.Ctx) error {

	var (
		userColl = database.GetCollection("users")
		data     userAuth.SignupReqDto
		otpColl  = database.GetCollection("otp")
	)
	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(userAuth.SignupResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if email is not already used
	filter := bson.M{
		"email":    strings.ToLower(data.Email),
		"isDeleted": false,
	}

	exists, err := userColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(userAuth.SignupResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(userAuth.SignupResDto{
			Status:  false,
			Message: "Email is already in use",
		})
	}

	otp := utils.Generate6DigitOtp()

	smallEmail := strings.ToLower(data.Email)
	// Store the OTP in the forgotPassword collection
	otpData := entity.OtpEntity{
		Id:        primitive.NewObjectID(),
		Otp:       otp,
		Email:     smallEmail,
		CreatedAt: time.Now().UTC(),
	}

	_, err = otpColl.InsertOne(ctx, otpData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.SignupResDto{
			Status:  false,
			Message: "failed to store OTP in the database: " + err.Error(),
		})
	}
	// Send OTP to the provided email
	_, err = utils.SendEmail(data.Email, otp)
	if err != nil {
		return c.Status(500).JSON(userAuth.SignupResDto{
			Status:  false,
			Message: "failed to send OTP email: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(userAuth.SignupResDto{
		Status:  true,
		Message: "otp send successfully",
	})
}
