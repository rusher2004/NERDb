package null

import (
	"database/sql"
	"encoding/json"
)

type JSONNullFloat64 struct {
	sql.NullFloat64
}

func (f JSONNullFloat64) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return json.Marshal(f.Float64)
	}

	return []byte("null"), nil
}

func (f *JSONNullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		f.Float64, f.Valid = 0, false
		return nil
	}

	f.Valid = true
	return json.Unmarshal(data, &f.Float64)
}

type JSONNullInt32 struct {
	sql.NullInt32
}

func (i JSONNullInt32) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int32)
	}

	return []byte("null"), nil
}

func (i *JSONNullInt32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.Int32, i.Valid = 0, false
		return nil
	}

	i.Valid = true
	return json.Unmarshal(data, &i.Int32)
}

type JSONNullInt64 struct {
	sql.NullInt64
}

func (i JSONNullInt64) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int64)
	}

	return []byte("null"), nil
}

func (i *JSONNullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.Int64, i.Valid = 0, false
		return nil
	}

	i.Valid = true
	return json.Unmarshal(data, &i.Int64)
}

type JSONNullString struct {
	sql.NullString
}

func (s JSONNullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}

	return []byte("null"), nil
}

func (s *JSONNullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.String, s.Valid = "", false
		return nil
	}

	s.Valid = true
	return json.Unmarshal(data, &s.String)
}
