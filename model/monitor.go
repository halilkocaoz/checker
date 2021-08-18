package model

import (
	"net/http"
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
	m.Headers = nil
}

// gets and sets post form values from database
func (m *Monitor) SetPostValues() {
	m.PostValues = nil
}
