package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// RelationTuple defines the structure of a relationship tuple for Ory Keto
type RelationTuple struct {
	Namespace  string `json:"namespace"`
	Object     string `json:"object"`
	Relation   string `json:"relation,omitempty"`
	SubjectID  string `json:"subject_id,omitempty"`
	SubjectSet *struct {
		Namespace string `json:"namespace"`
		Object    string `json:"object"`
		Relation  string `json:"relation"`
	} `json:"subject_set,omitempty"`
}

// SeedKetoDataFromFile reads relation tuples from a JSON file and sends them to Ory Keto's Write API
func SeedKetoDataFromFile(filePath string, apiURL string) error {
	// Read the JSON file
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	// Parse the JSON data
	var data []RelationTuple
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Loop through each relation tuple and send to Keto API
	for _, relation := range data {
		err := createRelationship(apiURL, relation)
		if err != nil {
			return fmt.Errorf("error creating relationship: %v", err)
		}
	}

	fmt.Println("Data seeded successfully.")
	return nil
}

// createRelationship sends a single relation tuple to Ory Keto's Write API using PUT /admin/relation-tuples
func createRelationship(apiURL string, relation RelationTuple) error {
	// Convert the relation to JSON
	jsonData, err := json.Marshal(relation)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Make the request to Keto Admin API (using PUT)
	req, err := http.NewRequest(http.MethodPut, apiURL+"/admin/relation-tuples", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create relationship, status code: %d", resp.StatusCode)
	}

	return nil
}
