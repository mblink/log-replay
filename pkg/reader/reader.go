package reader

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// LogEntry is single parsed entry from the log file
type LogEntry struct {
	Time      time.Time
	Method    string
	URL       string
	Payload   string
	UserAgent string
}

// LogReader provides generic log parser interface
type LogReader interface {
	Read(line string) (*LogEntry, error)
}

func ParseRequest(requestString string) ([]string, error) {
	parsedRequest := strings.SplitN(requestString, " ", 3)

	if len(parsedRequest) != 3 {
		return parsedRequest, fmt.Errorf("ERROR while parsing string: %s", requestString)
	}

	return parsedRequest, nil
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
