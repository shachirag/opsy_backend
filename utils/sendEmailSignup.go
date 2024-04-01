package utils

import (
	"context"
	"opsy_backend/database"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func SendSignupEmail(to string, link string) (*ses.SendEmailOutput, error) {
	var (
		subject2    = aws.String("Verify Your Opsy Account Email")
		senderEmail = os.Getenv("SENDER_EMAIL")
		charSet     = aws.String("UTF-8")
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
					Data:    SignupUserOtpEmailBodyHtml(link),
					Charset: charSet,
				},
				Text: &types.Content{
					Data:    SignupUserOtpEmailBodyText(link),
					Charset: charSet,
				},
			},
			Subject: &types.Content{
				Data:    subject2,
				Charset: charSet,
			},
		},
		Source: aws.String(senderEmail),
	}

	return sesClient.SendEmail(context.Background(), input)
}
