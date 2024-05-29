package logging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.JSONFormatter{}

	// Open the file for writing
	file, err := os.OpenFile("logging.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		// Set the logger to write logs to the file
		logger.Out = file
	} else {
		logger.Error("Failed to open logging.txt file for writing")
	}
}

// LogError logs an error event
func LogError(event, message string, err error) {
	fields := logrus.Fields{
		"event":   event,
		"message": message,
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	logger.WithFields(fields).Error("An error occurred")
}

// LogUserCreation logs user creation event
func LogUserCreation(name, game string) {
	logger.WithFields(logrus.Fields{
		"event":     "UserCreation",
		"operation": "creating user",
		"createdAt": time.Now(),
		"Name":      name,
		"Game":      game,
	}).Info("User was created")
}

// LogUserLogin logs user login event
func LogUserLogin(email string) {
	logger.WithFields(logrus.Fields{
		"event": "UserLogin",
		"email": email,
		"time":  time.Now(),
	}).Info("User logged in")
}

// LogValidation logs user validation event
func LogValidation(email string) {
	logger.WithFields(logrus.Fields{
		"event": "Validation",
		"user":  email,
		"time":  time.Now(),
	}).Info("User validated")
}
