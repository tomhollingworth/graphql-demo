// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package domain

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type DateTimeRange struct {
	Max time.Time `json:"max"`
	Min time.Time `json:"min"`
}

type EquipmentProperty struct {
	ID      string     `json:"id"`
	History []*History `json:"history,omitempty"`
}

func (EquipmentProperty) IsEntity() {}

type FilterHistory struct {
	Timestamp  *DateTimeRange `json:"timestamp,omitempty"`
	PropertyID *string        `json:"propertyID,omitempty"`
}

type History struct {
	EquipmentProperty *EquipmentProperty `json:"equipmentProperty"`
	Timestamp         time.Time          `json:"timestamp"`
	Value             string             `json:"value"`
	Datatype          DataType           `json:"datatype"`
}

type NewHistory struct {
	PropertyID string    `json:"propertyID"`
	Timestamp  time.Time `json:"timestamp"`
	Value      string    `json:"value"`
	Datatype   DataType  `json:"datatype"`
}

type DataType string

const (
	DataTypeString DataType = "STRING"
	DataTypeInt    DataType = "INT"
	DataTypeFloat  DataType = "FLOAT"
)

var AllDataType = []DataType{
	DataTypeString,
	DataTypeInt,
	DataTypeFloat,
}

func (e DataType) IsValid() bool {
	switch e {
	case DataTypeString, DataTypeInt, DataTypeFloat:
		return true
	}
	return false
}

func (e DataType) String() string {
	return string(e)
}

func (e *DataType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DataType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DataType", str)
	}
	return nil
}

func (e DataType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
