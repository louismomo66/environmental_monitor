package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func RespondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding data to JSON: %v", err)
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}
type CustomError struct {
	Message string `json:"message"`
}

func SetError(err error, msg string) CustomError {
	if err != nil {
		return CustomError{
			Message: msg + ": " + err.Error(),
		}
	}
	return CustomError{
		Message: msg,
	}
}

func SendEmail(to, subject, body string) error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
		return  err
	}
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	// Setup the SMTP configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully to", to)
	return nil
}