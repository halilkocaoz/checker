package model

import (
	"database/sql"
	"fmt"
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

	Deleted bool

	Headers    []KVPair
	PostValues []KVPair
}

func (m *Monitor) Process() {
	for {
		m.reGet()
		if m.Deleted { // if it's marked as deleted, free it
			m = nil
			return
		}

		m.setHeaders()
		if m.Method == "POST" {
			m.setPostValues()
		}

		resp, err := m.doRequest()
		if err != nil {
			// todo log
			fmt.Println(err)
		} else {
			resp.Body.Close()
			go m.writeResponseToDatabase(resp)
			if resp.StatusCode < 400 {
				//todo if it's not succes, push notifier service
			}
			fmt.Printf("%d: %v\n", resp.StatusCode, m)
		}

		time.Sleep(time.Duration(m.IntervalMs) * time.Millisecond)
	}
}

func (m *Monitor) writeResponseToDatabase(resp *http.Response) {
	insertResponseStatement := `INSERT INTO "Responses" ("MonitorID", "StatusCode") VALUES ($1, $2);`
	db, _ := storage.UpsMoDBConn()
	defer db.Close()
	db.Exec(insertResponseStatement, m.ID, resp.StatusCode)
}

// reget gets the monitor from database and pass new values.
// if it's deleted or its region is changed, marks as deleted.
func (m *Monitor) reGet() {
	db, _ := storage.UpsMoDBConn()
	defer db.Close()

	if err := db.QueryRow(`SELECT "Host",
	"Method",
	"IntervalMs",
	"TimeoutMs"
	FROM "Monitors"
	WHERE ("DeletedAt" IS NULL AND "Region" = $1 AND "ID" = $2)`, m.Region, m.ID).
		Scan(&m.Host,
			&m.Method,
			&m.IntervalMs,
			&m.TimeoutMs); err != nil {
		if err == sql.ErrNoRows {
			m.Deleted = true
		}
	}
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
