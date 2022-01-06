package haproxy

import (
	"fmt"
	"strings"
	"time"

	"github.com/Gonzih/log-replay/pkg/reader"
)

const (
	haProxyTsLayout = "2/Jan/2006:15:04:05.000"
)

// HaproxyReader implements reader.LogReader intefrace
type HaproxyReader struct {}

func parseHaproxyTime(timeLocal string) time.Time {
	t, err := time.Parse(haProxyTsLayout, timeLocal)

	reader.Must(err)

	return t
}

func parseStringInto(s string, entry *reader.LogEntry) error {
	dateStartI := strings.LastIndex(s, "[") + 1
	dateEndI := strings.LastIndex(s, "]")

	if dateStartI > dateEndI || dateStartI > len(s) || dateEndI > len(s) {
		return fmt.Errorf("Issue with date indexes, start: %d, end: %d, len: %d", dateStartI, dateEndI, len(s))
	}

	requestStartI := strings.Index(s, `"`) + 1
	requestEndI := len(s) - 1

	if requestStartI > requestEndI || requestStartI > len(s) || requestEndI > len(s) {
		return fmt.Errorf("Issue with request indexes, start: %d, end: %d, len: %d", requestStartI, requestEndI, len(s))
	}

	dateString := s[dateStartI:dateEndI]
	requestString := s[requestStartI:requestEndI]

	parsedRequest, err := reader.ParseRequest(requestString)

	if err != nil {
		return err
	}

	entry.Method = parsedRequest[0]
	entry.URL = parsedRequest[1]
	entry.Time = parseHaproxyTime(dateString)
	entry.UserAgent = ""

	return nil
}

// NewReader creates new reader for a haproxy log format
func NewReader() reader.LogReader {
	var reader HaproxyReader
	return &reader
}

func (r *HaproxyReader) Read(line string) (*reader.LogEntry, error) {
	var entry reader.LogEntry
	reader.Must(parseStringInto(line, &entry))
	return &entry, nil
}
