package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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
	from := "your-email@example.com"
	password := "your-email-password"

	// Setup the SMTP configuration
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}