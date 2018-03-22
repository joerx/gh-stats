package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type csvFormatter struct {
	w              *csv.Writer
	headers        []string
	fields         []string
	headersWritten bool
}

// NewCSVFormatter outputs as CSV
func NewCSVFormatter(w io.Writer, fields []string, headers []string) Formatter {
	if headers == nil {
		headers = fields
	}
	return &csvFormatter{w: csv.NewWriter(w), fields: fields, headers: headers}
}

func mayDeref(v reflect.Value) reflect.Value {
	k := v.Kind()
	if k.String() == "ptr" {
		return v.Elem()
	}
	return v
}

func toStructValue(d interface{}) (reflect.Value, bool) {
	v := mayDeref(reflect.ValueOf(d))
	if v.Kind().String() == "struct" {
		return v, true
	}
	return v, false
}

// Format renders data as CSV
func (f *csvFormatter) Write(data interface{}) error {
	structVal, isStruct := toStructValue(data)

	if !isStruct {
		return errors.New("Not a struct")
	}

	if !f.headersWritten {
		f.w.Write(f.headers)
		f.headersWritten = true
	}

	row := make([]string, len(f.fields))

	for i, n := range f.fields {
		fieldVal := mayDeref(structVal.FieldByName(n))
		if !fieldVal.IsValid() {
			row[i] = "n/a"
		} else {
			row[i] = fmt.Sprint(fieldVal)
		}
	}
	return f.w.Write(row)
}

func (f *csvFormatter) Flush() {
	f.w.Flush()
}
