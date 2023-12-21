package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch user By ID
// @Description Fetch user By ID
// @Tags user authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "business user ID"
// @Produce json
// @Success 200 {object} userAuth.GetUserResDto
// @Router /user/get-info/{id} [get]
func FetchUserById(c *fiber.Ctx) error {

	var user entity.UserEntity

	// Get the business user ID from the URL parameter
	userId := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return c.Status(400).JSON(userAuth.GetUserResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})

	}

	userColl := database.GetCollection("users")

	err = userColl.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(userAuth.GetUserResDto{
				Status:  false,
				Message: " user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.GetUserResDto{
			Status:  false,
			Message: "Failed to fetch user from MongoDB: " + err.Error(),
		})
	}

	// Convert business user to the response model (customerRes)
	userRes := userAuth.UserData{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(userAuth.GetUserResDto{
		Status:  true,
		Message: "user data retrieved successfully",
		User:    userRes,
	})
}
