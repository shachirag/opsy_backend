package logEntry
import (
    "encoding/json"
    "fmt"
    "opsy_backend/database"
    logEntry "opsy_backend/dto/users/logEntry"
    "opsy_backend/entity"
    "time"
 
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)
 
// @Summary fetch all required data
// @Tags logEntry
// @Description fetch all required data
// @Produce json
// @Param date query string true "date (YYYY-MM-DD)"
// @Success 200 {object} logEntry.CatgoriesResDto
// @Router /user/fetch-all-data [get]
func Months(c *fiber.Ctx) error {
    var (
        logEntryColl = database.GetCollection("logEntry")
    )
   monthStr:=c.Query("month")
   yearStr:=c.Query("year")
   date, err := time.Parse("2006-01-02", yearStr+"-"+monthStr+"-01")
   if err != nil {
	   return c.Status(fiber.StatusBadRequest).JSON(logEntry.CatgoriesResDto{
		   Status:  false,
		   Message: "Invalid date format: " + err.Error(),
	   })
   }

    // Adjust the format of todayStart to match the createdAt format in MongoDB
 
    filter := bson.M{
        "isDeleted": false,
		"createdAt": bson.M{"$gte": date, "$lte":getLastDateOfMonth(date)},
    }
    sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})
    // Log the filter value
    filterJSON, _ := json.Marshal(filter)
    fmt.Println(string(filterJSON))
 
    // Fetch data based on filters
    cursor, err := logEntryColl.Find(ctx, filter, sortOptions)
    if err != nil {
        fmt.Println("Error fetching data:", err.Error()) // Log the error
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
            fmt.Println("Error decoding data:", err.Error()) // Log the decoding error
            return c.Status(fiber.StatusInternalServerError).JSON(logEntry.CatgoriesResDto{
                Status:  false,
                Message: "Failed to decode data: " + err.Error(),
            })
        }
        logEntryData = append(logEntryData, logEntry.CategoriesRes{
            Id:   logEntryEntity.Id,
            Type: logEntryEntity.Type,
            When: logEntry.When{
                Date: logEntryEntity.When.Date,
                Time: logEntryEntity.When.Time,
            },
            PainLevel: logEntryEntity.PainLevel,
            Feel:      logEntryEntity.Feel,
            CreatedAt: logEntryEntity.CreatedAt,
            UpdatedAt: logEntryEntity.UpdatedAt,
        })
    }
 
    // Check if the result set is empty
    if len(logEntryData) == 0 {
        return c.Status(fiber.StatusOK).JSON(logEntry.CatgoriesResDto{
            Status:  false,
            Message: "No data found.",
        })
    }
 
    return c.Status(fiber.StatusOK).JSON(logEntry.CatgoriesResDto{
        Status:  true,
        Message: "Successfully fetched the data.",
        Data:    logEntryData,
    })
}


func getLastDateOfMonth(t time.Time) time.Time {
	// Get the first day of the next month
	firstDayOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())

	// Subtract one day to get the last day of the current month
	lastDayOfMonth := firstDayOfNextMonth.Add(-time.Second)

	return lastDayOfMonth
}