package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Verify OTP for signup
// @Description Verify the entered 6 digit OTP
// @Tags user authorization
// @Accept application/json
// @Param user body userAuth.VerifyOtpSignupReqDto true "Verify 6 digit OTP and insert data into database"
// @Produce json
// @Success 200 {object} userAuth.VerifyOtpSignupResDto
// @Router /user/verify-otp-for-signup [post]
func VerifyOtpForSignup(c *fiber.Ctx) error {
	var (
		otpColl  = database.GetCollection("otp")
		userColl = database.GetCollection("users")
		data     userAuth.VerifyOtpSignupReqDto
		otpData  entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Entered OTP is required",
		})
	}
	smallEmail := strings.ToLower(data.Email)
	filter := bson.M{"email": smallEmail, "isDeleted": false}
	// Find the user with email address from client
	err = otpColl.FindOne(ctx, filter, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&otpData)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(userAuth.VerifyOtpSignupResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}

		return c.Status(500).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Internal server error, while getting the user: " + err.Error(),
		})
	}

	// Compare the entered OTP with the OTP from the database
	if data.EnteredOTP != otpData.Otp {

		return c.Status(400).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Invalid Crdentials" + err.Error(),
		})
	}

	// Check if email is not already used
	filter = bson.M{
		"email": strings.ToLower(data.Email),
	}

	exists, err := userColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Email is already in use",
		})
	}

	id := primitive.NewObjectID()

	// Now that OTP is verified, proceed to insert the data into the database
	user := entity.UserEntity{
		Id:        id,
		Name:      data.Name,
		Email:     smallEmail,
		IsDeleted: false,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = userColl.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Failed to insert user data: " + err.Error(),
		})
	}

	// create auth token
	_secret := os.Getenv("JWT_SECRET_KEY")
	// _token_exp := os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT")
	// t, err := utils.CreateToken(user, _secret)
	month := (time.Hour * 24) * 30
	// Create the JWT claims, which includes the user ID and expiry time
	claims := jtoken.MapClaims{
		"Id":    user.Id,
		"email": user.Email,
		"role":  "user",
		"exp":   time.Now().Add(month * 6).Unix(),
	}
	// Create token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	_token, err := token.SignedString([]byte(_secret))
	if err != nil {
		return c.Status(400).JSON(userAuth.VerifyOtpSignupResDto{
			Status:  false,
			Message: "Token is not valid" + err.Error(),
		})
	}

	return c.Status(200).JSON(userAuth.VerifyOtpSignupResDto{
		Status:  true,
		Message: "OTP verified successfully and user data inserted",
		Token:   _token,
		Data: userAuth.UserResDto{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
