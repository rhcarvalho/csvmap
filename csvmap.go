// Package csvmap is a companion to the standard library package encoding/csv.
//
// Unlike encoding/csv, it represents records as maps of column names to
// values, providing direct access to fields by column name and making it easy
// to filter and reorder columns in CSV files.
package csvmap

import (
	"encoding/csv"
)

// A Reader wraps a csv.Reader to read records as maps of column names to
// values, instead of lists of values. It assumes the first record of the CSV
// input is a header with column names.
type Reader struct {
	r      *csv.Reader
	header []string
}

// NewReader returns a new Reader that reads from r.
func NewReader(r *csv.Reader) (*Reader, error) {
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	return &Reader{
		r:      r,
		header: header,
	}, nil
}

// Header returns the CSV header, a list of column names.
func (r *Reader) Header() []string {
	return r.header
}

// HasColumn reports whether the CSV header contains a column with the given
// name.
func (r *Reader) HasColumn(name string) bool {
	for _, col := range r.header {
		if col == name {
			return true
		}
	}
	return false
}

// Read reads one record, a map of column names to values, from r.
func (r *Reader) Read() (map[string]string, error) {
	record, err := r.r.Read()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, min(len(r.header), len(record)))
	for i := range record {
		if i >= len(r.header) {
			break
		}
		out[r.header[i]] = record[i]
	}
	return out, nil
}

// A Writer wraps a csv.Writer to write records as maps of column names to
// values, instead of lists of values.
type Writer struct {
	w      *csv.Writer
	header []string
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w *csv.Writer, header []string) *Writer {
	return &Writer{
		w:      w,
		header: header,
	}
}

// WriteHeader writes the CSV header to w along with any necessary quoting.
func (w *Writer) WriteHeader() error {
	return w.w.Write(w.header)
}

// Writer writes a single CSV record to w along with any necessary quoting. A
// record is a map of column names to values. Only columns present in the
// Writer's header are written, and in the order they appear in the header.
func (w *Writer) Write(record map[string]string) error {
	s := make([]string, len(w.header))
	for i, name := range w.header {
		s[i] = record[name]
	}
	return w.w.Write(s)
}

// WriteAll writes multiple CSV records to w using Write and then calls w.Flush.
func (w *Writer) WriteAll(records []map[string]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	w.w.Flush()
	return w.w.Error()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
