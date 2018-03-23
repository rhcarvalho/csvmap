package csvmap_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rhcarvalho/csvmap"
)

func ExampleReader() {
	// An io.Reader with example CSV content. In a typical case, use os.Open
	// to open a CSV file.
	s := strings.NewReader(`name,age
John,8
Jane,12
James,23
`)

	// Create a new Reader.
	// To customize how the CSV content is to be parsed, change the
	// csv.Reader before calling csvmap.NewReader.
	r, err := csvmap.NewReader(csv.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}

	// Inspect header.
	fmt.Println(r.Header())
	fmt.Println(r.HasColumn("age"), r.HasColumn("height"))

	// Read and print all records.
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name: %s, Age: %s\n", record["name"], record["age"])
	}
	// Output:
	// [name age]
	// true false
	// Name: John, Age: 8
	// Name: Jane, Age: 12
	// Name: James, Age: 23
}

func ExampleWriter() {
	// Some records that could have been read using a csvmap.Reader.
	records := []map[string]string{
		map[string]string{
			"name":  "John",
			"email": "john@example.com",
			"age":   "8",
		},
		map[string]string{
			"name":  "Marie",
			"email": "marie@example.com",
			"age":   "6",
		},
	}

	// Choose what columns to output and the order.
	header := []string{"email", "name"}

	// Create a new Writer.
	// To customize how the CSV content is to be formatted, change the
	// csv.Writer before calling csvmap.NewWriter.
	w := csvmap.NewWriter(csv.NewWriter(os.Stdout), header)

	err := w.WriteHeader()
	if err != nil {
		log.Fatal(err)
	}

	err = w.WriteAll(records)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// email,name
	// john@example.com,John
	// marie@example.com,Marie
}
