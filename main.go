package main

import "opsy_backend/app"

// @title Opsy Backend
// @version 0.0.1
// @description Opsy Backend Backend (in GoLang)
// @contact.name Priyanshu Rai
// @license.name MIT
// @host nssdmvtpj2.us-east-1.awsapprunner.com
// @BasePath /
func main() {
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}