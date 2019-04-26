package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// InfluxDb - db params. used as a db conector
type InfluxDB struct {
	host     string
	port     string
	database string
	username string
	password string
}

// Close - does nothing, made for overall consistency
func (db *InfluxDB) Close() error {
	return nil
}

// QueryURL - urlencodes query string and populates database query url with db params
func (db *InfluxDB) QueryURL(query string) (string, error) {
	q, err := url.Parse(query)
	if err != nil {
		return "", err
	}
	query = q.String()
	return fmt.Sprintf("http://%s:%s/query?db=%s&u=%s&p=%s&q=%s", db.host, db.port, db.database, db.username, db.password, query), nil
}

// Select - makes an arbitrary select query to InfluxDB
// returns string
func (db *InfluxDB) Select(query string, format string) (string, error) {
	if !strings.HasPrefix(strings.ToLower(query), "select") {
		return "", errors.New("Must be a 'select' query")
	}

	formats := [...]string{"csv", "json"}
	var validFormat bool
	for _, f := range formats {
		if format == f {
			validFormat = true
			break
		}
	}
	if !validFormat {
		format = "csv"
	}
	u, err := db.QueryURL(query)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Accept", fmt.Sprintf("application/%s", format))
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(buf.String())
	}

	return buf.String(), nil
}

// NewInfluxDbConnector - returns InfluxDB populated with db params
func NewInfluxDbConnector() *InfluxDB {
	return &InfluxDB{
		host:     "host",
		port:     "port",
		database: "db",
		username: "user",
		password: "pass",
	}
}

func NewCustomInfluxConnector(stageConfig map[string]string) *InfluxDB {
	return &InfluxDB{
		host:     stageConfig["host"],
		port:     stageConfig["port"],
		database: stageConfig["database"],
		username: stageConfig["username"],
		password: stageConfig["password"],
	}
}
