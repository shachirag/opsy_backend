package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func SignupUserOtpEmailBodyText(otp string) *string {
	s := fmt.Sprintf(`## OTP for Signup

Please use the 6-digit OTP below for signup.

<strong style="font-size: 130%%">%v</strong>

If you didn’t request this, you can ignore this email.

Thanks,
***OPSY***`, otp)

	return aws.String(s)
}

func SignupUserOtpEmailBodyHtml(otp string) *string {
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
				<th colspan="2" style="text-align: left; padding: 15px 15px 15px 0; background-color: #fff; color: white;">
					<a href="#"><img style="width: 170px; height: auto;" src="[App_Url]/logo/opsy_logo.png" alt=""></a>
				</th>
			</tr>
			<tr>
				<td style="border: none; padding-left:0;">
				</td>
			</tr>
			<tr>
				<td colspan="2" style="padding: 15px; padding-left:0;">
					<p>Welcome to Opsy, your all-in-one app for managing your physical and mental health! To get started, please use the following One-Time Password (OTP) to verify your email and complete
					the registration process:
					</p>
					<p><B>OTP: [Your OTP Code]</B></p>
					<p>Once you've verified your email, you'll have access to Opsy's powerful features, including:
					</p>
					<p>
					<ul>
						<li style="list-style: decimal;"><b>Health Tracking:</b> Opsy simplifies the process of logging and monitoring your physical and mental health data.
						Whether you're tracking your pain level or any other health-related input, 
						Opsy provides an easy-to-use platform for monitoring and analyzing your overall well-being. </li>
						<li style="list-style: decimal;"><b>Appointments:</b>
							Keep track of your upcoming appointments effortlessly. Opsy allows you to input and manage your
							appointments, ensuring you never miss an important date.</li>
						<li style="list-style: decimal;">
							<b>Graphical Insights:</b>
							Gain valuable insights into your health trends with our intuitive weekly, monthly, and yearly
							graphs.
						</li>
					</ul>
					</p>
					<p>To complete your registration:
					<ul>
						<li style="list-style: decimal;">Open the Opsy app.</li>
						<li style="list-style: decimal;">Enter the email address you used for registration.</li>
						<li style="list-style: decimal;">Use the provided OTP to verify your email.</li>
						<li style="list-style: decimal;">Set up your personalized profile and start exploring Opsy's
							features.</li>
					</ul>
					</p>
					<p>If you have any questions or need assistance, feel free to reach out to our support team at
						<a style="color:#87CEEB ; text-decoration: underline;" href="">support@opsyapp.com.</a>
					</p>
					<p>Thank you for choosing Opsy to be your personal health tracker.</p>
					<p>Best Regards,<br> <b>The Opsy Team</b></p>
				</td>
			</tr>
			<tr>
				<td colspan="2" style="text-align: center; padding: 15px; background-color: #f4f4f4;">
					<p>This is a notification email from the Opsy Team. Please do not reply to this email.​</p>
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
