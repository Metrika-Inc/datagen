package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type datafile struct {
	filename string
	theFile  *os.File
	reader   *bufio.Reader
}

func newDatafile(filename string) *datafile {
	return &datafile{
		filename: filename,
	}
}

func (f *datafile) open() {
	file, err := os.Open(f.filename)
	if err != nil {
		panic("Could not open data file.")
	}

	f.theFile = file
	f.reader = bufio.NewReader(file)

	// discard the opening bracket
	f.reader.ReadString('\n')
}

func (f *datafile) close() {
	f.theFile.Close()
}

func (f *datafile) getChunk() ([]byte, error) {
	if f.theFile == nil {
		return nil, fmt.Errorf("invalid file handle")
	}

	outText := make([]byte, 0, 1000)
	// scan the next entry.
	openCurlies := 0
	first := true

	for (openCurlies > 0) || first {

		currentByte, err := f.reader.ReadByte()
		if err != nil {
			return nil, err
		}

		switch currentByte {
		case '{':
			first = false
			openCurlies++
		case '}':
			openCurlies--
		}

		outText = append(outText, currentByte)
	}

	// skip the comma.
	f.reader.ReadByte()
	return outText, nil
}

func main() {
	filename := flag.String("f", "ledger.json", "the filename of the JSON data file.")
	recordCount := flag.Int("r", 149, "the number of records to return each time the /data endpoint is called.")

	flag.Parse()

	// File I/O stuff.
	f := newDatafile(*filename)
	f.open()
	defer f.close()

	// Web serving stuff.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/data", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		response := make([]byte, 0, 1000)
		response = append(response, "[\n"...)
		for i := 0; i < *recordCount; i++ {
			chunk, err := f.getChunk()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			response = append(response, chunk...)
			response = append(response, ',')
		}
		response = response[:len(response)-1]
		response = append(response, "\n]\n"...)

		if len(response) == 4 {
			w.WriteHeader(204)
		}

		w.Write(response)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
