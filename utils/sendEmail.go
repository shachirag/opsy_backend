package utils

import (
	"context"
	"opsy_backend/database"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)



func SendEmail(to string, link string) (*ses.SendEmailOutput, error) {
	var (
		senderEmail = os.Getenv("SENDER_EMAIL")
		charSet = aws.String("UTF-8")
		subject = aws.String("otp for user signup")
	)
	sesClient := database.GetSesClient()

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				// to,
				"opsy@yopmail.com",
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    SignupUserOtpEmailBodyHtml(link),
					Charset: charSet,
				},
				Text: &types.Content{
					Data:    SignupUserOtpEmailBodyText(link),
					Charset: charSet,
				},
			},
			Subject: &types.Content{
				Data:    subject,
				Charset: charSet,
			},
		},
		Source:  aws.String(senderEmail),
	}

	return sesClient.SendEmail(context.Background(), input)
}
