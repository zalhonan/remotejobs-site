package entity

import "time"

type JobRaw struct {
	ID             int64
	Content        string
	Title          string
	SourceLink     string
	MainTechnology string
	ContentPure    string
	DatePosted     time.Time
	DateParsed     time.Time
}
