package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

var currentBarangay string

func postAnnouncement() {
	_, err := loadCurrentUserDetails()
	if err != nil {
		fmt.Printf("Error loading user details: %v\n", err)
		return
	}

	fmt.Println("Announcement: ")
	reader := bufio.NewReader(os.Stdin)
	announcement, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("error reading input: %v\n", err)
		return
	}
	announcement = strings.TrimSpace(announcement)

	// Load the barangay data
	file, err := os.OpenFile("barangays.csv", os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	var records [][]string

	// Read through the CSV file
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return
		}

		// Check if this is the barangay of the logged-in admin
		if record[0] == currentBarangay {
			record[2] = announcement // Update the announcement field
		}
		records = append(records, record)
	}

	// Clear the file and write the modified data back to it
	file.Truncate(0)
	file.Seek(0, io.SeekStart)
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		fmt.Printf("error writing to file: %v\n", err)
		return
	}
	writer.Flush()

	fmt.Println("Announcement posted successfully!")
}

func updateStatus() {
	scanner := bufio.NewScanner(os.Stdin)

	// Valid status options
	validStatuses := map[string]bool{
		"Safe":     true,
		"Monitor":  true,
		"Alert":    true,
		"Evacuate": true,
	}

	var newStatus string
	for {
		fmt.Print("Enter new status (Safe, Monitor, Alert, Evacuate): ")
		scanner.Scan()
		newStatus := scanner.Text()

		// Check if the entered status is valid
		if _, isValid := validStatuses[newStatus]; isValid {
			break
		} else {
			fmt.Println("Invalid status entered. Please enter one of the following: " +
				"Safe, Monitor, Alert, Evacuate.")
		}
	}

	updateCSV("barangays.csv", currentBarangay, newStatus, "")

	fmt.Println("Barangay status updated successfully!")
}

func updateCSV(filename, barangayName, newStatus, newAnnouncement string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) > 1 && fields[0] == barangayName {
			fields[1] = newStatus
			fields[2] = newAnnouncement
			line = strings.Join(fields, ",")
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	/* to write the updated lines back to the file */
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}

type User struct {
	Email       string
	Password    string
	AccountType string
	Barangay    string
}

func loadCurrentUserDetails() ([]User, error) {
	file, err := os.Open("users.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var users []User

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		// Assuming the CSV has the following columns:
		// Email, Password, AccountType, Barangay
		if len(record) >= 4 {
			user := User{
				Email:       record[0],
				Password:    record[1],
				AccountType: record[2],
				Barangay:    record[3],
			}
			if record[0] == currentUserEmail {
				users = append(users, user)
				//fmt.Println("List of Users:")
				for _, user := range users {
					//fmt.Printf("Email: %s, Password: %s, Account Type: %s, Barangay: %s\n",
					//	user.Email, user.Password, user.AccountType, user.Barangay)
					currentBarangay = user.Barangay
					//fmt.Println(current_barangay)
				}
			}
		}
	}
	return users, nil
}

func removeAnnouncement() {
	_, err := loadCurrentUserDetails()
	if err != nil {
		fmt.Printf("Error loading user details: %v\n", err)
		return
	}

	// Load the barangay data
	file, err := os.OpenFile("barangays.csv", os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	var records [][]string
	// Read through the CSV file
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return
		}

		// Check if this is the barangay of the logged-in admin
		if record[0] == currentBarangay {
			if record[2] == "" {
				fmt.Println("The announcement is empty. No deletion needed.")
				return
			}
			record[2] = "" // Update the announcement field
		}
		records = append(records, record)
	}

	// Clear the file and write the modified data back to it
	file.Truncate(0)
	file.Seek(0, io.SeekStart)
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		fmt.Printf("error writing to file: %v\n", err)
		return
	}
	writer.Flush()

	fmt.Println("Announcement removed successfully!")
}
