package domain

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type DateTime time.Time

func MarshalDateTime(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return nil
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(t.Format(time.RFC3339Nano)))
	})
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		return time.Parse(time.RFC3339Nano, tmpStr)
	}
	return time.Time{}, errors.New("time should be RFC3339Nano formatted string")
}
