package log

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ErrorTracker logs errors to an external service
type ErrorTracker interface {
	Name() string
	TrackError(err error, tags map[string]string)
}

type singleton struct {
	tracker ErrorTracker
}

var instance *singleton
var once sync.Once

func getInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

// RegisterErrorTracker registers an error tracker which is called when errors are raised
func RegisterErrorTracker(errTracker ErrorTracker) {
	inst := getInstance()
	inst.tracker = errTracker
	Info("Error tracker registered", map[string]string{
		"tracker": errTracker.Name(),
	})
}

func print(level string, message string, tags map[string]string) {
	om := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"message":   message,
		"tags":      tags,
	}
	output, err := json.Marshal(om)
	if err != nil {
		fmt.Printf("Failed to marshal log json: %s", err.Error())
		return
	}
	fmt.Println(string(output))
}

// Info logs messages and optional metdata to the console
func Info(message string, tags map[string]string) {
	print("INFO", message, tags)
}

// Error logs an error and optional metadata to console and to any registered ErrorTracker
func Error(err error, tags map[string]string) {
	print("ERROR", err.Error(), tags)

	inst := getInstance()
	if inst.tracker != nil {
		inst.tracker.TrackError(err, tags)
	}
}
