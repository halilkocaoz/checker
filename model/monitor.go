package model

import (
	"net/http"
	"net/url"
	"time"

	"github.com/halilkocaoz/upsmo-checker/storage"
)

// UpsMo monitor
type Monitor struct {
	ID         string
	Host       string
	Method     string
	Region     string
	IntervalMs int
	TimeoutMs  int
	CreatedAt  string

	Headers    []KVPair
	PostValues []KVPair
}

// does http request and return result
func (monitor *Monitor) doRequest() (*http.Response, error) {
	client := http.Client{}
	client.Timeout = time.Duration(time.Millisecond * time.Duration(monitor.TimeoutMs))
	hostUrl, err := url.Parse(monitor.Host)
	if err != nil {
		return nil, err
	}
	request := new(http.Request)
	request.URL = hostUrl
	request.Method = monitor.Method
	request.Header = make(http.Header)

	for _, header := range monitor.Headers {
		request.Header.Add(header.Key, header.Value)
	}
	if monitor.Method == "POST" && len(monitor.PostValues) > 0 {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		form := url.Values{}
		for _, postform := range monitor.PostValues {
			form.Add(postform.Key, postform.Value)
		}
		request.PostForm = form
	}

	return client.Do(request)
}

// gets and sets header values from database
func (m *Monitor) setHeaders() {
	headers := make([]KVPair, 0)
	db, _ := storage.UpsMoDBConn()
	defer db.Close()

	headerRows, _ := db.Query(`SELECT 
	"Key", 
	"Value"
	FROM "Headers" 
	WHERE ("DeletedAt" IS NULL AND "MonitorID" = $1) 
	ORDER BY "CreatedAt"`, m.ID)
	defer headerRows.Close()

	for headerRows.Next() {
		header := new(KVPair)
		_ = headerRows.Scan(&header.Key, &header.Value)
		headers = append(headers, *header)
	}
	m.Headers = headers
}

// gets and sets post form values from database
func (m *Monitor) setPostValues() {
	postForms := make([]KVPair, 0)
	db, _ := storage.UpsMoDBConn()
	defer db.Close()

	postFormRows, _ := db.Query(`SELECT
	"Key", 
	"Value"
	FROM "PostForms" 
	WHERE ("DeletedAt" IS NULL AND "MonitorID" = $1) 
	ORDER BY "CreatedAt"`, m.ID)
	defer postFormRows.Close()

	for postFormRows.Next() {
		postForm := new(KVPair)
		_ = postFormRows.Scan(&postForm.Key, &postForm.Value)
		postForms = append(postForms, *postForm)
	}
	m.PostValues = postForms
}
