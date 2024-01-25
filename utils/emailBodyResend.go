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
	s := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml">
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Reset your password and login</title>
			<!--[if mso]><style type="text/css">body, table, td, a { font-family: Arial, Helvetica, sans-serif !important; }</style><![endif]-->
		</head>
		<body style="font-family: Helvetica, Arial, sans-serif; margin: 0px; padding: 0px; background-color: #ffffff;">
			<table role="presentation"
				style="width: 100%%; border-collapse: collapse; border: 0px; border-spacing: 0px; font-family: Arial, Helvetica, sans-serif; background-color: rgb(255, 255, 255);">
				<tbody>
					<tr>
						<td align="center" style="padding: 1rem 2rem; vertical-align: top; width: 100%%;">
							<table role="presentation" style="max-width: 600px; border-collapse: collapse; border: 0px; border-spacing: 0px; text-align: left;">
								<tbody>
									<tr>
										<td style="padding: 40px 0px 0px;">
											<div style="text-align: center;">
												<div style="padding-bottom: 20px;"><img src="[App_Url]/logo/opsy_logo.png" alt="" style="width: 100px;"></div>
											</div>
											<div style="padding: 20px; background-color: rgb(237, 237, 237);">
												<div style="color: rgb(19, 17, 18); text-align: center;">
													<h3>Otp for reset password</h3>
													<p style="padding-bottom: 16px">Please use the 6-digit OTP below for for reset password.</p>
													<center style="padding-top: 8px; padding-bottom: 20px;"><strong style="font-size: 170%%;">[OTP]</strong></center>
	
													<p style="padding-bottom: 16px">If you didn’t request this, you can ignore this email.</p>
	
													<p style="padding-bottom: 16px">Thanks,<br><em><strong>Careville</strong></em></p>
												</div>
											</div>
											<div style="padding-top: 20px; color: rgb(254, 98, 98); text-align: center;">
												<p style="padding-bottom: 16px">Made with ♥ in Nigeria</p>
											</div>
										</td>
									</tr>
								</tbody>
							</table>
						</td>
					</tr>
				</tbody>
			</table>
		</body>
	</html>`
	s = strings.ReplaceAll(s, "[Your OTP Code]", otp)
	s = strings.ReplaceAll(s, "[App_Url]", os.Getenv("S3_BUCKET_URL"))
	str := aws.String(s)
	return str
}
