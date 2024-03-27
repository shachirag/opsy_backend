package utils

import (
	"context"
	"opsy_backend/database"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func SendEmail(to string, name string, link string) (*ses.SendEmailOutput, error) {
	var (
		senderEmail = os.Getenv("SENDER_EMAIL")
		charSet     = aws.String("UTF-8")
		subject     = aws.String("Reset Your Opsy Password")
	)
	sesClient := database.GetSesClient()

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    ForgotUserOtpEmailBodyHtml(name, link),
					Charset: charSet,
				},
				Text: &types.Content{
					Data:    ForgotSignupUserOtpEmailBodyText(link),
					Charset: charSet,
				},
			},
			Subject: &types.Content{
				Data:    subject,
				Charset: charSet,
			},
		},
		Source: aws.String(senderEmail),
	}

	return sesClient.SendEmail(context.Background(), input)
}
