package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonFormatter struct {
	w io.Writer
}

// Format renders data as JSON
func (f *jsonFormatter) Write(data interface{}) error {
	b, _ := json.Marshal(data)
	fmt.Fprint(f.w, string(b))
	return nil
}

func (f *jsonFormatter) Flush() {
	// noop
}

// NewJSONFormatter outputs as JSON
func NewJSONFormatter(w io.Writer) Formatter {
	return &jsonFormatter{w: w}
}
