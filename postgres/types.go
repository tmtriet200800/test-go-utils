package pkgPostgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	pkgError "github.com/tmtriet200800/test-go-utils/errors"
)

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct{ sql.NullInt64 }

// MarshalJSON for NullInt64
func (v NullInt64) MarshalJSON() ([]byte, error) {
    if v.Valid {
        return json.Marshal(v.Int64)
    } else {
        return json.Marshal(nil)
    }
}

func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
        return err
	}
	
	if x != nil {
        v.Valid = true
        v.Int64 = *x
    } else {
        v.Valid = false
	}
    return nil
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct{ sql.NullBool }

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}

	jsn, err := json.Marshal(nb.Bool)
	if err != nil {
		return jsn, fmt.Errorf("Postgres could not marshal NullBool: %w", err)
	}

	return jsn, nil
}

// UnmarshalJSON for NullBool
func (nb NullBool) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &nb.Bool); err != nil {
		return pkgError.Wrap(fmt.Errorf("Postgres NullBool unmarshal error: %w", err))
	}

	nb.Valid = true

	return nil
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct{ sql.NullFloat64 }

// MarshalJSON for NullFloat64
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}

	jsn, err := json.Marshal(nf.Float64)
	if err != nil {
		return jsn, fmt.Errorf("Postgres could not marshal NullFloat64: %w", err)
	}

	return jsn, nil
}

// UnmarshalJSON for NullFloat64
func (nf NullFloat64) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &nf.Float64); err != nil {
		return pkgError.Wrap(fmt.Errorf("Postgres NullFloat64 unmarshal error: %w", err))
	}

	nf.Valid = true

	return nil
}

// NullString is an alias for sql.NullString data type
type NullString struct{ sql.NullString }

func (v NullString) MarshalJSON() ([]byte, error) {
    if v.Valid {
        return json.Marshal(v.String)
    } else {
        return json.Marshal(nil)
    }
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
        return err
	}
	
	if x != nil {
        v.Valid = true
        v.String = *x
    } else {
        v.Valid = false
	}
    return nil
}

// NullTime is an alias for mysql.NullTime data type
type NullTime struct{ sql.NullTime }

func (v NullTime) MarshalJSON() ([]byte, error) {
    if v.Valid {
        return json.Marshal(v.Time)
    } else {
        return json.Marshal(nil)
    }
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
        return err
	}
	
	if x != nil {
        v.Valid = true
        v.Time = *x
    } else {
        v.Valid = false
	}
    return nil
}
