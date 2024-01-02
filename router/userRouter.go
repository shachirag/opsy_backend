package router

import (
	"opsy_backend/handlers"
	"opsy_backend/handlers/users/logEntry"
	userAuthenticate "opsy_backend/handlers/users/userAuthentication"

	"github.com/gofiber/fiber/v2"
)

func UsersSetupsRoutes(app *fiber.App) {

	app.Static("/", ".puplic")

	/* ---------- HEALTH ---------- */
	app.Get("/health", handlers.HealthCheck)

	/* ---------- Protected Routes ----- */
	// secret := os.Getenv("JWT_SECRET_KEY")
	// jwt := middlewares.NewAuthMiddleware(secret)

	//user authentication

	user := app.Group("/user")

	/* ----- admin authentication -----*/
	user.Post("/login", userAuthenticate.LoginUser)
	user.Post("/signup", userAuthenticate.SignupUser)
	user.Post("/verify-otp-for-signup", userAuthenticate.VerifyOtpForSignup)
	user.Post("/forgot-password", userAuthenticate.ForgotPassword)
	user.Post("/verify-otp", userAuthenticate.VerifyOtp)
	user.Put("/reset-password", userAuthenticate.ResetPasswordAfterOtp)
	user.Post("/resend-otp", userAuthenticate.ResendOTP)
	user.Put("/update-user-data/:id", userAuthenticate.UpdateUser)
	user.Get("/get-info/:id",userAuthenticate.FetchUserById)
	user.Put("/change-password/:id",userAuthenticate.ChangeUserPassword)
	user.Get("/get-misc-data/:id",userAuthenticate.FetchAllMiscData)
	//log Entries
	user.Post("/create-log-entry",logEntry.CreateLogEntry)
}
