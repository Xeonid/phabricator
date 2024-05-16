package old

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	. "phabricator/bugz"
	"strconv"
	"strings"
)

func getIDS() {
	// Specify the path to the directory containing JSON files
	directoryPath := "bugsJson/"

	fmt.Println("Extracting IDS from bugs")

	// Get a list of all JSON files in the directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Slice to store extracted IDs
	var allIDs []int

	// Iterate over each file in the directory
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			// Construct the full path to the JSON file
			filePath := filepath.Join(directoryPath, file.Name())

			// Read the JSON file
			jsonData, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading JSON file %s: %v\n", filePath, err)
				continue
			}

			// Process the JSON data
			var bug Bug
			err = json.Unmarshal(jsonData, &bug)
			if err != nil {
				fmt.Printf("Error decoding JSON from file %s: %v\n", filePath, err)
				continue
			}

			// Extract 'id' fields
			ids := extractIDs(bug)

			// Append to the slice of all IDs
			allIDs = append(allIDs, ids...)
		}
	}

	// Remove duplicates from the slice
	allIDs = removeDuplicates(allIDs)

	// Write the unique IDs to a file
	err = writeTooFile("old", "uniqueIDs1.txt", allIDs)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("Unique IDs written to uniqueIDs1.txt")
}

func extractIDs(bug Bug) []int {
	var ids []int

	// Extracting 'id' fields
	ids = appendIfNotExists(ids, bug.AssignedToDetail.ID)

	// Extract 'id' field from CCDetail
	for _, ccUser := range bug.CCDetail {
		ids = appendIfNotExists(ids, ccUser.ID)
	}

	// Extract 'id' field from CreatorDetail
	ids = appendIfNotExists(ids, bug.CreatorDetail.ID)

	return ids
}

func appendIfNotExists(slice []int, id int) []int {
	for _, existingID := range slice {
		if existingID == id {
			return slice
		}
	}
	return append(slice, id)
}

func writeTooFile(directory, filename string, ids []int) error {
	// Convert IDs to string
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, strconv.Itoa(id))
	}

	// Join IDs with a newline character
	data := []byte(strings.Join(idsStr, "\n"))

	// Create the full file path
	fullPath := filepath.Join(directory, filename)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}

	// Write the data to the specified file
	return os.WriteFile(fullPath, data, 0644)
}

func removeDuplicates(slice []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range slice {
		if encountered[slice[v]] == true {
			// Already seen, skip
		} else {
			// Record as an encountered element
			encountered[slice[v]] = true
			// Append to result
			result = append(result, slice[v])
		}
	}
	return result
}
