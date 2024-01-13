package userAuthenticate

import (
	"context"
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"
	"os"
	"strings"
	"time"

	jtoken "github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

// @Summary Login User
// @Description Login User
// @Tags user authorization
// @Accept application/json
// @Param user body userAuth.LoginReqDto true "Login User"
// @Produce json
// Success 200 {object} userAuth.LoginResDto
// @Router /user/login [post]
func LoginUser(c *fiber.Ctx) error {

	var (
		userColl = database.GetCollection("users")
		data     userAuth.LoginReqDto
		user     entity.UserEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(userAuth.LoginResDto{
			Status:  false,
			Message: "failed to parse data" + err.Error(),
		})
	}

	smallEmail := strings.ToLower(data.Email)
	err = userColl.FindOne(ctx, bson.M{"email": smallEmail, "isDeleted": false}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(500).JSON(userAuth.LoginResDto{
				Status:  false,
				Message: "Email does not exist",
			})
		}
		return c.Status(500).JSON(userAuth.LoginResDto{
			Status:  false,
			Message: "internal server error" + err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(strings.TrimSpace(data.Password)))
	if err != nil {
		return c.Status(500).JSON(userAuth.LoginResDto{
			Status:  false,
			Message: "invalid credentials",
		})
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	month := (time.Hour * 24) * 30

	claims := jtoken.MapClaims{
		"Id":    user.Id,
		"email": user.Email,
		"role":  "user",
		"exp":   time.Now().Add(month * 6).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	generatedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(400).JSON(userAuth.LoginResDto{
			Status:  false,
			Message: "token is not valid" + err.Error(),
		})
	}

	return c.Status(200).JSON(userAuth.LoginResDto{
		Status:  true,
		Message: "Successfully logged in",
		Data: userAuth.UserResDto{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: generatedToken,
	})
}
