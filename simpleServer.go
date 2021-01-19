package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	// Since we are not interested in the contents of the file
	// we are only concerned with the error result
	_, readFileErr := ioutil.ReadFile("tmp/ok")
	if readFileErr != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "tmp/ok not found")
	} else {
		fmt.Fprintf(w, "OK")
	}
	return
}

func imgHandler(w http.ResponseWriter, r *http.Request) {

	// Move the request into print friendly(ish) json
	jsonHeader, jsonMarshallErr := json.Marshal(r.Header)
	if jsonMarshallErr != nil {
		simpleLog("Error in request header: %s", jsonMarshallErr)
	}

	// Log request
	simpleLog("img request [%s]", jsonHeader)

	// Open up the 1x1 gif
	imgFile, fileOpenErr := os.Open("tracker.gif")
	if fileOpenErr != nil {
		simpleLog("Error opening tracker.gif: %s", fileOpenErr)
		return
	}
	defer imgFile.Close()

	// Read file into slice to get the content type and size
	header := make([]byte, 512)
	imgFile.Read(header)
	fileType := http.DetectContentType(header)

	stat, fileStatErr := imgFile.Stat()
	if fileStatErr != nil {
		simpleLog("Error getting tracker.gif file statistics")
	}
	size := strconv.FormatInt(stat.Size(), 10)

	// Populate the response header
	w.Header().Set("Content-Type", fileType)
	w.Header().Set("Content-Disposition", "attachment; filename=tracker.gif")
	w.Header().Set("Content-Length", size)

	// Write the file with the response writer
	imgFile.Seek(0, 0)
	io.Copy(w, imgFile)
	return
}

func main() {
	simpleLog("Starting simple tracking webserver")
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/img", imgHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func simpleLog(format string, a ...interface{}) {
	curTime := time.Now().UTC()
	logMessage := fmt.Sprintf(format, a...)
	fmt.Printf("[%s] %s\n", curTime.String(), logMessage)
	return
}
