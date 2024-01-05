package logEntry
 
import (
    "context"
    "opsy_backend/database"
    logEntry "opsy_backend/dto/users/logEntry"
    "opsy_backend/entity"
    "time"
 
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
)
 
var ctx = context.Background()
 
// @Summary CreateLogEntry
// @Description  CreateLogEntry
// @Tags logEntry
// @Param user body logEntry.LogEntryReqDto true "CreateLogEntry for user"
// @Produce json
// @Success 200 {object} logEntry.GetLogEntryResDto
// @Router /user/create-log-entry [post]
func CreateLogEntry(c *fiber.Ctx) error {
 
    var (
        logEntryColl = database.GetCollection("logEntry")
        data         logEntry.LogEntryReqDto
    )
    // Parsing the request body
    err := c.BodyParser(&data)
    if err != nil {
        return c.Status(500).JSON(logEntry.GetLogEntryResDto{
            Status:  false,
            Message: err.Error(),
        })
    }
   
    parsedDate, err := time.Parse("2006-01-02", data.Date)
    if err != nil {
        return c.Status(500).JSON(logEntry.GetLogEntryResDto{
            Status:  false,
            Message: err.Error(),
        })
    }
 
    parsedTime, err := time.Parse("15:04", data.Time)
    if err != nil {
        return c.Status(500).JSON(logEntry.GetLogEntryResDto{
            Status:  false,
            Message: err.Error(),
        })
    }
 
    mergedDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
        parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)
 
    logEntryData := entity.LogEntryEntity{
        Id:          primitive.NewObjectID(),
        Type:        data.Type,
        IsDeleted:   false,
        Feel:        data.Feel,
        Notes:       data.Notes,
        Ways:        data.Ways,
        When:        mergedDateTime,
        PainLevel:   data.PainLevel,
        WhatItIsFor: data.WhatItIsFor,
        Alert:       data.Alert,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }
 
    _, err = logEntryColl.InsertOne(ctx, logEntryData)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(logEntry.GetLogEntryResDto{
            Status:  false,
            Message: "failed to store Log Entry Database " + err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(logEntry.GetLogEntryResDto{
        Status:  true,
        Message: "Data inserted successfully",
    })
}
