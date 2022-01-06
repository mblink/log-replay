package solr

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Gonzih/log-replay/pkg/reader"
)

const (
	solrProxyTsLayout = "2006-01-02 15:04:05.000"
)

// SolrReader implements reader.LogReader intefrace
type SolrReader struct {}

func parseSolrTime(timeLocal string) time.Time {
	t, err := time.Parse(solrProxyTsLayout, timeLocal)

	reader.Must(err)

	return t
}

func parseSolrPayload(params string) (string, error) {
	r := regexp.MustCompile(`{(.+)?}`)
	matches := r.FindStringSubmatch(params)
	if len(matches) != 2 {
		return "", fmt.Errorf("Unable to parse solr payload.")
	}
	return matches[1], nil
}

func parseSolrInto(s string, entry *reader.LogEntry) error {
	if len(s) < 23 {
		return fmt.Errorf("This log line does not seem to contain a valid timestamp.")
	}
	dateString := strings.Replace(s[0:23], ",", ".", -1)
	stringParts := strings.Split(s, " ")

	var requestParts [2]string
	for _, part := range stringParts {
		if strings.HasPrefix(part, "path") {
			requestParts[0] = part
		} else if strings.HasPrefix(part, "params") {
			requestParts[1] = part
		}
	}
	//Default solr requests to post to go around query length for GET requests.
	payload, err := parseSolrPayload(requestParts[1])
	if err != nil {
		return err
	}

	path := strings.SplitAfterN(requestParts[0], "=", 2)

	entry.Method = "POST"
	entry.URL = path[1]
	entry.Time = parseSolrTime(dateString)
	entry.Payload = payload
	return nil
}

// NewReader creates new reader for a solr log format
func NewReader() reader.LogReader {
	var reader SolrReader
	return &reader
}

func (r *SolrReader) Read(line string) (*reader.LogEntry, error) {
	var entry reader.LogEntry
	reader.Must(parseSolrInto(line, &entry))
	return &entry, nil
}
