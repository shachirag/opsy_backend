package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// var ctx = context.Background()

// @Summary Fetch All Categories
// @Description Fetch All Categories
// @Tags expertise for users
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} expertise.CatgoriesResDto
// @Router /user/get-categories [get]
func FetchAllMiscData(c *fiber.Ctx) error {
	// Get the misc collection
	miscColl := database.GetCollection("misc")

	// Fetch the categories data-
	var miscdata userAuth.CategoriesRes
	categoryFilter := bson.M{"_id": "phycicalHealth"}
	cursor, err := miscColl.Find(ctx, categoryFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.CatgoriesResDto{
			Status:  false,
			Message: "Failed to fetch categories data: " + err.Error(),
		})
	}

	for cursor.Next(ctx) {
		err := cursor.Decode(&miscdata)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(userAuth.CatgoriesResDto{
				Status:  false,
				Message: "Failed to decode  data: " + err.Error(),
			})
		}
	}

	

	// Prepare the response
	response := userAuth.CatgoriesResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data:    miscdata,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
