package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func ResendUserOtpEmailBodyText(otp string) *string {
	s := fmt.Sprintf(`## Otp for User reset password

Please use the 6-digit OTP below for reset password.

<strong style="font-size: 130%%">%v</strong>

If you didn’t request this, you can ignore this email.

Thanks,
***Careville***`, otp)

	return aws.String(s)
}

func ResendUserOtpEmailBodyHtml(otp string) *string {
	s := `<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title></title>
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				margin: 0;
				padding: 0;
			}
	
			table {
				width: 100%;
				max-width: 600px;
				margin: 0 auto;
			}
	
			th,
			td {
				padding: 10px;
				text-align: left;
				border-bottom: 1px solid #ddd;
			}
	
			th {
				background-color: #4CAF50;
				color: white;
			}
		</style>
	</head>
	
	<body>
		<table>
			<tr>
				<th colspan="2" style="text-align: left; padding:15px 15px 15px 0; background-color: #fff; color: white;">
					<a href="#"><img style="width: 170px; height: auto;" src="[App_Url]/logo/opsy_logo.png" alt=""></a>
				</th>
			</tr>
			<tr>
				<td style="border: none; padding-left:0;">
					<h2>Resend OTP for Opsy Account</h2>
				</td>
			</tr>
			<tr>
				<td colspan="2" style="padding: 15px; padding-left:0;">
					<p>Dear User,</p>
					<p>We noticed that you're in the process of securing your Opsy account, and we're here to assist you.
						use the following One-Time Password (OTP) to verify your email:
					</p>
					<p><B>OTP: [Your OTP Code]</B></p>
					<p>No need to worry about complicated steps—simply use this code in the Opsy app to proceed.
					</p>
					<p>
						If you have any questions or encounter any issues, our support team is ready to help. Reach out to
						us at <a href="">support@opsyapp.com</a>, and we'll get back to you promptly.
					</p>
					<p>Thank you for choosing Opsy to support your health and well-being journey.
					</p>
					<p>Best Regards,<br> <b>The Opsy Team</b></p>
				</td>
			</tr>
			<tr>
				<td colspan="2" style="text-align: center; padding: 15px; background-color: #f4f4f4;">
					<p>This is a notification email from <b>The Opsy Team </b>. Please do not reply to this
						email.</p>
				</td>
			</tr>
		</table>
	</body>
	
	</html>`
	s = strings.ReplaceAll(s, "[Your OTP Code]", otp)
	s = strings.ReplaceAll(s, "[App_Url]", os.Getenv("S3_BUCKET_URL"))
	str := aws.String(s)
	return str
}
