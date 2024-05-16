package old

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	. "phabricator/bugz"
	"strconv"
	"strings"
)

func GetUsersJson() {
	apiURL := "https://bugs.freebsd.org/bugzilla/rest/user"
	token, _ := getToken()
	// Create a 'users' folder if it doesn't exist
	err := os.Mkdir("usersJson", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating 'users' folder: %v\n", err)
		return
	}

	//Extract unique IDS from bugs
	getIDS()

	// Open the file containing IDs
	idFile, err := os.Open("old/uniqueIDs1.txt")
	if err != nil {
		fmt.Printf("Error opening id file: %v\n", err)
		return
	}

	// Create a scanner to read IDs from the file
	scanner := bufio.NewScanner(idFile)

	// Read IDs in a loop until there are no more left
	for scanner.Scan() {
		id := scanner.Text()

		// Create query parameters
		params := url.Values{}
		params.Set("ids", id)
		params.Set("token", token)

		// Construct the full URL with query parameters
		fullURL := apiURL + "?" + params.Encode()

		// Add print statements for debugging
		fmt.Println("Fetching data for ID:", id)

		// Make a GET request to the API
		response, err := http.Get(fullURL)
		if err != nil {
			fmt.Printf("Error making GET request to %s: %v\n", fullURL, err)
			continue
		}

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading response body from %s: %v\n", fullURL, err)
			response.Body.Close()
			continue
		}

		response.Body.Close()

		// Check if the response indicates an error
		if isError(body) || isEmptyUsersArray(body) {
			//fmt.Println("Skipping due to error in response")
			continue
		}

		// Unmarshal JSON data into User struct
		var userResponse map[string][]User
		err = json.Unmarshal(body, &userResponse)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			continue
		}

		// Iterate over users and write to individual files
		for _, user := range userResponse["users"] {
			// Convert the id to an integer
			idInt, _ := strconv.Atoi(id)

			filename := filepath.Join("usersJson", fmt.Sprintf("user_%d.json", idInt))
			err := writeToFile(filename, user)
			if err != nil {
				fmt.Printf("Error writing to file %s: %v\n", filename, err)
				return
			}
		}

	}
	fmt.Println("User data written to individual files.")
}
func writeToFile(filename string, user User) error {
	// Use json.MarshalIndent to preserve the original formatting
	indentedData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	// Write the data to the specified file
	return os.WriteFile(filename, indentedData, 0644)
}

func isError(response []byte) bool {
	return strings.Contains(string(response), "error\":true") || len(response) == 0
}

func isEmptyUsersArray(response []byte) bool {
	return strings.Contains(string(response), "\"users\":[]")
}
