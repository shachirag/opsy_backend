package logEntry

import (
	"opsy_backend/database"
	logEntry "opsy_backend/dto/users/logEntry"
	"opsy_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary fetch all requied data
// @Tags logEntry
// @Description fetch all requied data
// @Produce json
// @Success 200 {object} logEntry.CatgoriesResDto
// @Router /customer/fetch-all-data [get]
func FetchAllData(c *fiber.Ctx) error {
	var (
		logEntryColl = database.GetCollection("logEntry")
	)
	filter := bson.M{"isDeleted": false}
	// Define the sort options (descending order of updatedAt)
	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	// Fetch data based on filters
	cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "Failed to fetch data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var logEntryData []logEntry.CategoriesRes
	for cursor.Next(ctx) {
		var logEntryEntity entity.LogEntryEntity
		if err := cursor.Decode(&logEntryEntity); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CatgoriesResDto{
				Status:  false,
				Message: "Failed to decode data: " + err.Error(),
			})
		}
		logEntryData = append(logEntryData, logEntry.CategoriesRes{
			Id:        logEntryEntity.Id,
			Type:      logEntryEntity.Type,
			When:      logEntry.When{
				Date: logEntryEntity.When.Date,
				Time: logEntryEntity.When.Time,
			},
			PainLevel: logEntryEntity.PainLevel,
			Feel:      logEntryEntity.Feel,
			CreatedAt: logEntryEntity.CreatedAt,
			UpdatedAt: logEntryEntity.UpdatedAt,
		},
		)
	}

	// Check if is empty
	if len(logEntryData) == 0 {
		return c.Status(fiber.StatusOK).JSON(logEntry.CatgoriesResDto{
			Status:  false,
			Message: "No data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(logEntry.CatgoriesResDto{
		Status:  true,
		Message: "Successfully fetched the data.",
		Data:    logEntryData, // Include the fetched data
	})
}
