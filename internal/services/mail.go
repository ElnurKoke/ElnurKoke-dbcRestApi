package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"path/filepath"
)

type EmailService interface {
	SendEmail(user, password string, to []string, filename string, fileData interface{}) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s *emailService) SendEmail(user, password string, to []string, filename string, fileData interface{}) error {
	auth := smtp.PlainAuth("", user, password, "smtp.gmail.com")
	log.Println(user, "-", password, "=", to, "=", filename)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	textPart, err := writer.CreateFormField("text")
	if err != nil {
		return fmt.Errorf("failed to create email text part: %w", err)
	}
	_, _ = textPart.Write([]byte("Please find the attached file."))

	filePart, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("failed to create file part: %w", err)
	}

	if fileReader, ok := fileData.(multipart.File); ok {
		_, err = io.Copy(filePart, fileReader)
		if err != nil {
			return fmt.Errorf("failed to copy file data: %w", err)
		}
	} else {
		return fmt.Errorf("invalid file data")
	}

	writer.Close()

	msg := fmt.Sprintf(
		"To: %s\r\nSubject: File Delivery\r\nMIME-Version: 1.0\r\nContent-Type: %s\r\n\r\n%s",
		to[0], writer.FormDataContentType(), body.String(),
	)

	err = smtp.SendMail("smtp.gmail.com:587", auth, user, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
