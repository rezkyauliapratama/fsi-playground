package main

import (
	"flag"
	"fmt"
	"internal-iam-service/services"
	"os"
)

func main() {
	// Define flags for JSON file and API URL
	jsonFilePath := flag.String("file", "", "Path to the JSON file containing relation tuples")
	apiURL := flag.String("api-url", "http://localhost:4467/relation-tuples", "Ory Keto Write API URL")
	flag.Parse()

	// Validate that the JSON file path is provided
	if *jsonFilePath == "" {
		fmt.Println("Error: You must provide a valid path to a JSON file using --file flag.")
		os.Exit(1)
	}

	// Run the seeding process using the provided JSON file
	err := services.SeedKetoDataFromFile(*jsonFilePath, *apiURL)
	if err != nil {
		fmt.Println("Error seeding data:", err)
		os.Exit(1)
	}

	fmt.Println("Seeding completed successfully.")
}
