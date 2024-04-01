package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func ForgotSignupUserOtpEmailBodyText(otp string) *string {
	s := fmt.Sprintf(`## OTP for Reset password

Please use the 6-digit OTP below for Reset password.

<strong style="font-size: 130%%">%v</strong>

If you didn’t request this, you can ignore this email.

Thanks,
***OPSY***`, otp)

	return aws.String(s)
}

func ForgotUserOtpEmailBodyHtml(name string, otp string) *string {
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
				</td>
			</tr>
			<tr>
				<td colspan="2" style="padding: 15px; padding-left:0;">
					<p>Dear [User],</p>
					<p>We hope this email finds you well. It seems like you've requested to reset your password for your
						Opsy account. No worries – we're here to help you secure your account.
					</p>
					<p>To reset your password, please use the following One-Time Password (OTP) and follow the steps below:
					</p>
					<p><B>OTP: [Your OTP Code]</B></p>
					<p>Reset Your Password Steps:
					<ul>
						<li style="list-style: decimal;">Open the Opsy app.</li>
						<li style="list-style: decimal;">Navigate to the login screen.</li>
						<li style="list-style: decimal;">Select the "Forgot Password" option.</li>
						<li style="list-style: decimal;">Enter your email address.</li>
						<li style="list-style: decimal;">Use the provided OTP to verify your identity.</li>
						<li style="list-style: decimal;">Set up a new, secure password for your Opsy account.</li>
					</ul>
					</p>
					<p>
						Please note that this OTP is valid for a limited time to ensure the security of your account.
						</ul>
					</p>
					<p>
						If you did not request a password reset, or if you have any concerns, please contact our support
						team immediately at <a href="">support@opsyapp.com</a>. We take the security of your account
						seriously and are here
						to assist you.
					</p>
					<p>Thank you for trusting Opsy with your health journey. We're committed to providing a safe and secure
						environment for you to track your health and well-being.<br>
					</p>
					<p>Best Regards,<br> <b>The Opsy Team</b></p>
				</td>
			</tr>
			<tr>
				<td colspan="2" style="text-align: center; padding: 15px; background-color: #f4f4f4;">
					<p>This is a notification email from <b>The Opsy Team.</b> Please do not reply to this
						email.</p>
				</td>
			</tr>
		</table>
	</body>
	
	</html>`
	s = strings.ReplaceAll(s, "[Your OTP Code]", otp)
	s = strings.ReplaceAll(s, "[User]", name)
	s = strings.ReplaceAll(s, "[App_Url]", os.Getenv("S3_BUCKET_URL"))
	str := aws.String(s)
	return str
}
