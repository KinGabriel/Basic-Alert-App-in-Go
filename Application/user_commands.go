package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

var currentAlertStatus string
var currentAnnouncement string

type barangayAlert struct {
	barangay     string
	alertStatus  string
	announcement string
}

func loadCurrentBarangayDetails() ([]barangayAlert, error) {
	file, err := os.Open("barangays.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var brgys []barangayAlert

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		// Assuming the CSV has the following columns:
		// Barangay, Captain, Alert, Announcement
		if len(record) >= 4 {
			brgy := barangayAlert{
				barangay:     record[0],
				alertStatus:  record[1],
				announcement: record[2],
			}
			if record[0] == currentBarangay {
				brgys = append(brgys, brgy)
				//fmt.Println("List of Users:")
				for _, brgy := range brgys {
					currentAlertStatus = brgy.alertStatus
					currentAnnouncement = brgy.announcement
				}
			}
		}
	}
	return brgys, nil
}

func displayAnnouncement() {
	fmt.Printf("Announcement for %s : %s \n", currentBarangay, currentAnnouncement)
}

func displayAlertStatus() {
	fmt.Printf("Alert Status for %s : %s \n", currentBarangay, currentAlertStatus)
}

func displaySafetyProtocol() {
	// Open the file
	file, err := os.Open("safety_protocol.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read and print the file's content line by line
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Check for errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func searchBarangay() {
	var barangayChoice int
	validBarangays, err := loadBarangaysFromCsv()
	if err != nil {
		fmt.Printf("error loading barangays: %v", err)
	}
	displayBarangayChoices(validBarangays)
	fmt.Print("Enter the number of your Barangay: ") // ask to enter a barangay number
	fmt.Scanln(&barangayChoice)
	if barangayChoice > 0 && barangayChoice <= len(validBarangays) {
		barangay := validBarangays[barangayChoice-1]
		fmt.Printf("You have selected: %s\n", barangay)
		file, err := os.OpenFile("barangays.csv", os.O_RDONLY, os.ModePerm)
		if err != nil { // error handler if csv is empty
			fmt.Printf("error opening file: %v", err)
		}
		defer file.Close()
		reader := csv.NewReader(file)    // create a reader
		records, err := reader.ReadAll() // read all contents
		if err != nil {                  // error handler in the reading process
			fmt.Printf("error reading file: %v", err)
		}
		for _, record := range records { // loop the data in the csv
			if record[0] == barangay {
				fmt.Println(barangay)
				fmt.Println("Current status: " + record[1])
				fmt.Println("Current announcement: " + record[2])
				fmt.Println("Press enter to continue...")
				fmt.Scanln()
				break
			}
		}
	} else {
		fmt.Println("Invalid choice. Please select a valid Barangay.")
	}
}
