package zsql

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"
)

var nullBytes = []byte("null")

// A JSON safe version of sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullBytes, nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.Int64 = 0
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

type NullInt32 struct {
	sql.NullInt32
}

func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullBytes, nil
	}
	return json.Marshal(ni.Int32)
}

func (ni *NullInt32) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.Int32 = 0
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int32)
	ni.Valid = (err == nil)
	return err
}

type NullInt16 struct {
	sql.NullInt16
}

func (ni NullInt16) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullBytes, nil
	}
	return json.Marshal(ni.Int16)
}

func (ni *NullInt16) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.Int16 = 0
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int16)
	ni.Valid = (err == nil)
	return err
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (ni NullFloat64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullBytes, nil
	}
	return json.Marshal(ni.Float64)
}

func (ni *NullFloat64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.Float64 = 0
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Float64)
	ni.Valid = (err == nil)
	return err
}

type NullTime struct {
	sql.NullTime
}

func (ni NullTime) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullBytes, nil
	}
	return json.Marshal(ni.Time)
}

func (ni *NullTime) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.Time = time.Time{}
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Time)
	ni.Valid = (err == nil)
	return err
}

type NullString struct {
	sql.NullString
}

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return nullBytes, nil
	}
	return json.Marshal(ni.String)
}

func (ni *NullString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.String = ""
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.String)
	ni.Valid = (err == nil)
	return err
}
