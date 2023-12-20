package utils

import (
	"context"
	"opsy_backend/database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var (
	// senderEmail = os.Getenv("SENDER_EMAIL")
	senderEmail = "opsy@yopmail.com"
	charSet     = aws.String("UTF-8")
	sender      = aws.String(senderEmail)
	subject     = aws.String("otp for user signup")
)

func SendEmail(to string, link string) (*ses.SendEmailOutput, error) {
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
				Data:    subject,
				Charset: charSet,
			},
		},
		Source: sender,
	}

	return sesClient.SendEmail(context.Background(), input)
}
