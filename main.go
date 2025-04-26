package main

import (
	"encoding/csv"
	"flag" // Import the flag package
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gregdel/pushover"
)

type CsvLine struct {
	Day     int
	Month   int
	Year    int
	Subject string
	Event   string
	Notes   string
}

func extractCsvData(line []string) (c CsvLine, err error) {
	c = CsvLine{}
	day, err := strconv.Atoi(line[0])
	if err != nil {
		return c, fmt.Errorf("Error parsing day from CSV file, column value: %s", line[0])
	}

	month, err := strconv.Atoi(line[1])
	if err != nil {
		return c, fmt.Errorf("Error parsing month from CSV file, column value: %s", line[1])
	}

	year := 0
	if line[2] != "" {
		year, err = strconv.Atoi(line[2])
		if err != nil {
			return c, fmt.Errorf("Error parsing year from CSV file, column value: %s", line[2])
		}
	}
	data := CsvLine{
		Day:     day,
		Month:   month,
		Year:    year,
		Subject: line[3],
		Event:   line[4],
		Notes:   line[5],
	}

	return data, nil
}

func sendMessages(events []CsvLine) {
	poToken := os.Getenv("PUSHOVER_API_TOKEN")
	if poToken == "" {
		log.Fatal("Missing PUSHOVER_API_TOKEN")
	}

	poRecipient := os.Getenv("PUSHOVER_RECIPIENT_TOKEN")
	if poRecipient == "" {
		log.Fatal("Missing PUSHOVER_RECIPIENT_TOKEN")
	}

	app := pushover.New(poToken)
	recipient := pushover.NewRecipient(poRecipient)

	for _, e := range events {
		m := fmt.Sprintf("%s - %s\n", e.Subject, e.Event)
		message := pushover.NewMessage(m)
		log.Print(m)
		_, err := app.SendMessage(message, recipient)
		if err != nil {
			log.Panic(err)
		}
	}
}

func main() {
	// Define a string flag "file" with a default value "events.csv" and a help message.
	csvFilename := flag.String("file", "events.csv", "Path to the events CSV file")
	// Parse the command-line flags
	flag.Parse()

	log.Print("START")
	// Use the filename provided by the flag (dereference the pointer)
	f, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatalf("Couldn't open the csv file '%s': %v\n", *csvFilename, err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatalln("Error parsing CSV file", err)
	}

	var events []CsvLine

	for i, line := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		csvLine, err := extractCsvData(line)
		if err != nil {
			msg := fmt.Sprintf("Error on line %d", i+1)
			log.Fatalln(msg, err)
		}

		t := time.Now()
		d := t.Day()
		m := int(t.Month())

		if csvLine.Month == m && csvLine.Day == d {
			fmt.Println(csvLine.Subject)
			events = append(events, csvLine)
		}
	}

	log.Printf("Found %d events\n", len(events))

	if len(events) != 0 {
		sendMessages(events)
	}
	log.Print("END")
}
