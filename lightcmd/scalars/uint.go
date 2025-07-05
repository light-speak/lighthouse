package scalars

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalUint converts uint to string when sending to frontend
func MarshalUint(i uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		str := strconv.FormatUint(uint64(i), 10)
		_, _ = io.WriteString(w, `"`+str+`"`) // Wrap in quotes to ensure string output
	})
}

// UnmarshalUint converts string/number from frontend to uint
func UnmarshalUint(v any) (uint, error) {
	switch v := v.(type) {
	case string:
		u64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(u64), nil
	case int:
		if v < 0 {
			return 0, errors.New("cannot convert negative numbers to uint")
		}
		return uint(v), nil
	case int64:
		if v < 0 {
			return 0, errors.New("cannot convert negative numbers to uint")
		}
		return uint(v), nil
	case json.Number:
		u64, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(u64), nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}

// MarshalUint64 converts uint64 to string when sending to frontend
func MarshalUint64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		str := strconv.FormatUint(i, 10)
		_, _ = io.WriteString(w, `"`+str+`"`) // Wrap in quotes to ensure string output
	})
}

// UnmarshalUint64 converts string/number from frontend to uint64
func UnmarshalUint64(v any) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int:
		if v < 0 {
			return 0, errors.New("cannot convert negative numbers to uint64")
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, errors.New("cannot convert negative numbers to uint64")
		}
		return uint64(v), nil
	case json.Number:
		return strconv.ParseUint(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}

// MarshalUint32 converts uint32 to string when sending to frontend
func MarshalUint32(i uint32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		str := strconv.FormatUint(uint64(i), 10)
		_, _ = io.WriteString(w, `"`+str+`"`) // Wrap in quotes to ensure string output
	})
}

// UnmarshalUint32 converts string/number from frontend to uint32
func UnmarshalUint32(v any) (uint32, error) {
	switch v := v.(type) {
	case string:
		u64, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(u64), nil
	case int:
		if v < 0 {
			return 0, errors.New("cannot convert negative numbers to uint32")
		}
		return uint32(v), nil
	case int64:
		if v < 0 {
			return 0, errors.New("cannot convert negative numbers to uint32")
		}
		return uint32(v), nil
	case json.Number:
		u64, err := strconv.ParseUint(string(v), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(u64), nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}
