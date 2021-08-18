package model

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
