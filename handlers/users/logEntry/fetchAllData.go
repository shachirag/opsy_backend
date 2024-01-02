package logEntry

// import (
// 	"opsy_backend/database"
// 	"opsy_backend/dto/users/logEntry"
// 	"opsy_backend/entity"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // @Summary fetch all requied data
// // @Tags customer
// // @Description fetch all requied data
// // @Produce json
// // @Success 200 {object} logEntry.catagoriesRes
// // @Router /customer/fetch-all-data [get]
// func FetchAllData(c *fiber.Ctx) error {
//     var (
//         logEntryColl = database.GetCollection("logEntry")
//     )
//     filter := bson.M{"isDeleted": false}
//     // Define the sort options (descending order of updatedAt)
//     sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})
 
//     // Fetch salesAgent data based on filters
//     cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CategoriesRes{
//             Status:  false,
//             Message: "Failed to fetch data: " + err.Error(),
//         })
//     }
//     defer cursor.Close(ctx)
 
//     var logEntryData []logEntry.CategoriesRes
//     for cursor.Next(ctx) {
//         var logEntry entity.LogEntryEntity
//         if err := cursor.Decode(&logEntryData); err != nil {
//             return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CategoriesRes{
//                 Status:  false,
//                 Message: "Failed to decode sales agents data: " + err.Error(),
//             })
//         }
//         logEntryData = append(logEntryData,logEntry.CategoriesRes{
//             Id:          logEntry.Id,
// 			Type:        logEntry.Type,
// 			When:        logEntry.When,
// 			PainLvel:    logEntry.PainLevel,
// 			Feel:        logEntry.Feel,
// 			CreatedAt:   time.Now().UTC(),
// 			UpdatedAt:   time.Now().UTC(),
            
//             },
// 	)
//     }
// 	return c.Status(fiber.StatusOK).JSON(response)
// }


//     // Check if is empty
//     if len(logEntryData ) == 0 {
//         return c.Status(fiber.StatusOK).JSON(logEntry.CategoriesRes{
//             Status:  false,
//             Message: "No data found.",
//         })
//     }
 
//     return c.Status(fiber.StatusOK).JSON(logEntry.CategoriesRes{
//         Status:      true,
//         Message:     "Successfully fetched the data.",
//         logEntry: logEntryData, // Include the fetched data
//     })
// }