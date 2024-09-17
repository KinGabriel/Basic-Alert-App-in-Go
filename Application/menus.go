package main

import (
	"fmt"
	"os"
)

/*
Function which shows the start-up option of the program
*/
func LogInOrRegisterMenu() {
	for true { // loop incase of user errors
		fmt.Println("Welcome to the alert app of Baguio city")
		fmt.Println("====Options====")
		fmt.Println("[1]. Register")
		fmt.Println("[2]. Log In")
		fmt.Println("[3]. Exit")
		choiceLogOrReg()
	}
}

/*
Function which allows the user to choose if the user will log in,register or exit
*/
func choiceLogOrReg() {
	var choice int
	fmt.Print("Enter your choice: ") // Ask the user of its choice
	fmt.Scan(&choice)                // read the choice
	switch {
	case choice == 1:
		registerUser()
	case choice == 2: // To let the user register
		logInUser()
	case choice == 3:
		os.Exit(0)
	default: // error handler incase the user insert an invalid input
		fmt.Println("Invalid choice. Please enter a number between 1 and 3")
	}
}

/* Show the menus for the user */
func userMenu() {
	loadCurrentUserDetails()
	loadCurrentBarangayDetails()
	for true { // loop incase of user errors
		fmt.Println("1. See announcements")
		fmt.Println("2. See your barangay status")
		fmt.Println("3. Read about safety protocols")
		fmt.Println("4. Search other baranggays status and announcement")
		fmt.Println("5. Exit")
		userChoice()
	}
}

/* Let the user choose what he/she will use in the menus */
func userChoice() {
	var choice int
	fmt.Print("Enter your choice: ") // Ask the user of its choice
	fmt.Scan(&choice)                // read the choice
	switch {
	case choice == 1:
		displayAnnouncement()
	case choice == 2:
		displayAlertStatus()
	case choice == 3:
		displaySafetyProtocol()
	case choice == 4:
		searchBarangay()
	case choice == 5:
		LogInOrRegisterMenu()
	default: // error handler incase the user insert an invalid input
		fmt.Println("Invalid choice. Please enter a number between 1 and 5")
	}
}

func adminMenu() {
	for true {
		fmt.Println("1. Post an announcement")
		fmt.Println("2. Update status")
		fmt.Println("3. Remove Announcement")
		fmt.Println("4. Exit")
		adminChoice()
	}
}

func adminChoice() {
	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scan(&choice)
	switch {
	case choice == 1:
		postAnnouncement()
	case choice == 2:
		updateStatus()
	case choice == 3:
		removeAnnouncement()
	case choice == 4:
		LogInOrRegisterMenu()
	default: // error handler incase the user insert an invalid input
		fmt.Println("Invalid choice. Please enter a number between 1 and 4")
	}
}
