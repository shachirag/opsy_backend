package userAuthenticate

import (
	"opsy_backend/database"
	userAuth "opsy_backend/dto/users/userAuthentication"
	"opsy_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary Fetch All misc data
// @Description Fetch All misc data
// @Tags user authorization
// @Accept application/json
//
//	@Param Authorization header string true "Authentication header"
//
// @Produce json
// @Success 200 {object} userAuth.CatgoriesResDto
// @Router /user/get-misc-data [get]
func FetchAllMiscData(c *fiber.Ctx) error {
	// Get the misc collection
	miscColl := database.GetCollection("misc")

	var physicalHealth entity.MiscEntity
	physicalHealthFilter := bson.M{"_id": "physicalHealth"}
	err := miscColl.FindOne(ctx, physicalHealthFilter).Decode(&physicalHealth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.CatgoriesResDto{
			Status:  false,
			Message: "Failed to fetch physical health data: " + err.Error(),
		})
	}

	var mentalHealth entity.MiscEntity
	mentalHealthFilter := bson.M{"_id": "mentalHealth"}
	err = miscColl.FindOne(ctx, mentalHealthFilter).Decode(&mentalHealth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userAuth.CatgoriesResDto{
			Status:  false,
			Message: "Failed to fetch mental health data: " + err.Error(),
		})
	}

	// Map retrieved data to response structure
	response := userAuth.CatgoriesResDto{
		Status:  true,
		Message: "Misc data fetched successfully",
		Data: userAuth.CategoriesRes{
			PhysicalHealth: userAuth.PhysicalHealth{
				Popular: physicalHealth.Popular,
				Other:   physicalHealth.Other,
			},
			MentalHealth: userAuth.MentalHealth{
				Popular: mentalHealth.Popular,
				Other:   mentalHealth.Other,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
