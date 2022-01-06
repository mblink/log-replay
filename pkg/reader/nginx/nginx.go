package nginx

import (
	"time"

	"github.com/Gonzih/log-replay/pkg/reader"
	"github.com/satyrius/gonx"
)

const (
	nginxTimeLayout = "2/Jan/2006:15:04:05 -0700"
)

// NginxReader implements reader.LogReader intefrace
type NginxReader struct {
	GonxParser *gonx.Parser
}

func parseNginxTime(timeLocal string) time.Time {
	t, err := time.Parse(nginxTimeLayout, timeLocal)

	reader.Must(err)

	return t
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

	timeLocal, err := rec.Field("time_local")
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
	entry.Time = parseNginxTime(timeLocal)

	return &entry, nil
}
