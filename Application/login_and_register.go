package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

var currentUserEmail string

/* This method asks the user about the log-in credentials. */
func logInUser() {
	var email, password string
	for true {
		fmt.Print("Enter email: ") // ask for the user email
		fmt.Scan(&email)
		fmt.Print("Enter password: ") // ask for the user password
		fmt.Scan(&password)
		accountType, isValid := validateLogIn(email, password)
		if isValid { // validate the log in
			fmt.Println("Successful log in!")
			currentUserEmail = email
			switch strings.ToLower(accountType) {
			case "admin":
				adminMenu()
			case "user":
				userMenu()
			default:
				fmt.Println("Unknown account type")
			}
			return
		} else {
			fmt.Println("Incorrect email or password!")
			for true {
				var choice string
				fmt.Println("do you want to try again?[y/n]") // ask the user if he/she wants to try to log in again
				fmt.Scan(&choice)
				if choice == "y" || choice == "Y" {
					break // log in again
				} else if choice == "n" || choice == "N" {
					return // go back to the menu
				} else {
					fmt.Println("Inputs should only be y or n!") // in case of invalid inputs
				}
			}
		}
	}
}

/* This method checks the csv if the email and password matches. */
func validateLogIn(email, password string) (string, bool) {
	file, err := os.OpenFile("users.csv", os.O_RDONLY, os.ModePerm)
	if err != nil { // error handler if csv is empty
		fmt.Printf("error opening file: %v", err)
		return "", false
	}
	defer file.Close()
	reader := csv.NewReader(file) // create a reader

	records, err := reader.ReadAll() // read all contents
	if err != nil {                  // error handler in the reading process
		fmt.Printf("error reading file: %v", err)
	}

	for _, record := range records { // loop the data in the csv
		if record[0] == email && record[1] == password { // validate if there is a matching data based on the password and email
			return record[2], true // successful log in
		}
	}
	return "", false // unsuccessful log in
}

/* This method asks the user to fill out the fields inorder to register. */
func registerUser() {
	var email, password, retypePassword string
	var barangayChoice int
	fmt.Println("====Register====") // header
	allowedDomains := []string{"gmail", "yahoo", "outlook", "hotmail", "icloud"}
	validBarangays, err := loadBarangaysFromCsv()
	if err != nil {
		fmt.Printf("error loading barangays: %v", err)
		return
	}
	for {
		fmt.Print("Enter your email: ") // ask the user email
		fmt.Scanln(&email)
		if isEmailTaken(email) {
			fmt.Println("Email is taken!")
		} else if isValidEmail(email, allowedDomains) {
			break
		} else {
			fmt.Println("Invalid Email Address!")
		}
	}
	for {
		fmt.Print("Enter your password: ") // ask the user password
		fmt.Scanln(&password)
		if len(password) > 8 {
			break
		} else {
			fmt.Println("Invalid Password!, it must be at least 8 characters long")
		}
	}
	for {
		fmt.Print("Please retype your password: ") // ask to retype password
		fmt.Scanln(&retypePassword)
		if password != retypePassword {
			fmt.Println("Passwords don't match! please try again")
		} else {
			break
		}
	}
	// Barangay selection
	for {
		displayBarangayChoices(validBarangays)
		fmt.Print("Enter the number of your Barangay: ") // ask to enter a barangay number
		fmt.Scanln(&barangayChoice)

		if barangayChoice > 0 && barangayChoice <= len(validBarangays) {
			barangay := validBarangays[barangayChoice-1]
			fmt.Printf("You have selected: %s\n", barangay)
			registerToCsv(email, password, barangay)
			break
		} else {
			fmt.Println("Invalid choice. Please select a valid Barangay.")
		}
	}
}

/* This method checks the email if it is already registered in the csv file. */
func isEmailTaken(email string) bool {
	file, err := os.Open("users.csv")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error reading file: %v", err)
			return false
		}
		if record[0] == email {
			return true
		}
	}
	return false
}

/* This method checks the validity of the email. */
func isValidEmail(email string, allowedDomains []string) bool {
	if strings.Contains(email, "@") && strings.HasSuffix(email, ".com") {
		// Extract the domain part of the email
		parts := strings.Split(email, "@")
		if len(parts) < 2 {
			return false
		}
		domainPart := strings.Split(parts[1], ".")[0]
		for _, domain := range allowedDomains {
			if domain == domainPart {
				return true
			}
		}
	}
	return false
}

/* This method checks the barangay if it is valid based on the database.
func isValidBarangay(barangay string, validBarangays []string) bool {
	for _, brgy := range validBarangays {
		if strings.EqualFold(barangay, brgy) {
			return true
		}
	}
	return false
}
*/
/* This method loads the list of barangays from the database. */
func loadBarangaysFromCsv() ([]string, error) {
	file, err := os.Open("barangays.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}
	var barangays []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		barangayName := record[0]
		barangays = append(barangays, barangayName)
	}
	return barangays, nil
}

/* This method registers the user to the CSV file */
func registerToCsv(email, password, barangay string) bool {
	file, err := os.OpenFile("users.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil { // error handler if csv is empty
		fmt.Printf("error opening file: %v", err)
		return false
	}
	defer file.Close()
	writer := csv.NewWriter(file)                           // instantiate the writer
	register := []string{email, password, "user", barangay} // write the contents in the csv
	if err := writer.Write(register); err != nil {          // error handler in writing the contents
		fmt.Printf("error writing to file: %v", err)
		return false
	}
	writer.Flush()
	if err := file.Close(); err != nil { //error handler
		fmt.Printf("error closing file: %v", err)
	}
	fmt.Println("User registered successfully")
	return true
}

func displayBarangayChoices(barangays []string) {
	fmt.Println("Select your Barangay:")
	for i, barangay := range barangays {
		fmt.Printf("%d. %s\n", i+1, barangay)
	}
}
