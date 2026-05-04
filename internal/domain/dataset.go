package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("dataset not found")

type Dataset struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	FileName    string      `json:"fileName"`
	StoredName  string      `json:"storedName,omitempty"`
	ContentType string      `json:"contentType"`
	Size        int64       `json:"size"`
	Rows        int         `json:"rows"`
	Columns     []string    `json:"columns"`
	Profile     DataProfile `json:"profile"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Preview     []CSVRow    `json:"preview,omitempty"`
}

func (d Dataset) StorageName() string {
	if d.StoredName != "" {
		return d.StoredName
	}
	return d.FileName
}

type CSVRow map[string]string

type DataProfile struct {
	Numeric []NumericProfile `json:"numeric"`
}

type NumericProfile struct {
	Column string  `json:"column"`
	Count  int     `json:"count"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Avg    float64 `json:"avg"`
}
