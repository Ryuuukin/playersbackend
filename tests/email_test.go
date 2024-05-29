package tests

import (
	"testing" // Update with your actual package path

	"github.com/Ryuuukin/ap-assignment1/utils"
)

func TestSendMailSimple(t *testing.T) {
	// Define test data
	subject := "Test Subject"
	body := "Test Body"
	to := []string{"recipient@example.com"}

	// Call the function and check for errors
	err := utils.SendMailSimple(subject, body, to)
	if err != nil {
		t.Errorf("Error sending email: %v", err)
	}
}
