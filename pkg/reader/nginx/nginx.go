package nginx

import (
	"strconv"
	"strings"
	"time"

	"github.com/Gonzih/log-replay/pkg/reader"
	"github.com/satyrius/gonx"
)

// NginxReader implements reader.LogReader intefrace
type NginxReader struct {
	GonxParser *gonx.Parser
}

func parseNginxTime(msec string) time.Time {
	millis, err := strconv.ParseInt(strings.Replace(msec, ".", "", -1), 10, 64)
	reader.Must(err)

	return time.UnixMilli(millis)
}

// NewReader creates new reader for a haproxy log format using provided io.Reader
func NewReader(format string) reader.LogReader {
	var reader NginxReader
	reader.GonxParser = gonx.NewParser(format)
	return &reader
}

func (r *NginxReader) Read(line string) (*reader.LogEntry, error) {
	var entry reader.LogEntry

	rec, err := r.GonxParser.ParseString(line)
	if err != nil {
		return &entry, err
	}

	msec, err := rec.Field("msec")
	if err != nil {
		return &entry, err
	}

	requestString, err := rec.Field("request")
	if err != nil {
		return &entry, err
	}

	parsedRequest, err := reader.ParseRequest(requestString)
	if err != nil {
		return &entry, err
	}

	entry.Method = parsedRequest[0]
	entry.URL = parsedRequest[1]
	entry.Time = parseNginxTime(msec)

	return &entry, nil
}
