package model

import (
	"net/http"

	"github.com/halilkocaoz/upmo-checker/storage"
)

// upmo monitor
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
func (monitor *Monitor) DoRequest() (*http.Response, error) {
	return nil, nil
}

// gets and sets header values from database
func (m *Monitor) SetHeaders() {
	headers := make([]KVPair, 0)
	db, _ := storage.UpMoDBConnection()
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
func (m *Monitor) SetPostValues() {
	postForms := make([]KVPair, 0)
	db, _ := storage.UpMoDBConnection()
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
