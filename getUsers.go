package main

/*import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	. "phabricator/bugz"
)

// User represents user data
type User struct {
	ID       int     `json:"id"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	CanLogin bool    `json:"can_login"`
	RealName string  `json:"real_name"`
	Groups   []Group `json:"groups"`
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	apiURL := "https://bugs.freebsd.org/bugzilla/rest/user"

	// Generate an array of arbitrary size (e.g., 1000)
	ids := generateIDs(1000)

	for _, id := range ids {

		// Create query parameters
		params := url.Values{}
		params.Set("ids", strconv.Itoa(id))
		params.Set("token", "43929-vtmestP5W8")

		// Construct the full URL with query parameters
		fullURL := apiURL + "?" + params.Encode()

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
			fmt.Println("Skipping due to error in response")
			continue
		}

		// Write user details to individual files
		filename := filepath.Join("users", fmt.Sprintf("user_%d.txt", id))
		err = writeToFile(filename, body)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", filename, err)
			return
		}
	}

	fmt.Println("User data written to individual files.")
}

func writeToFile(filename string, data []byte) error {
	// Write the data to the specified file
	return os.WriteFile(filename, data, 0644)
}

func generateIDs(size int) []int {
	ids := make([]int, size)
	for i := 1; i <= size; i++ {
		ids[i-1] = i
	}
	return ids
}

func isError(response []byte) bool {
	return strings.Contains(string(response), "error\":true") || len(response) == 0
}

func isEmptyUsersArray(response []byte) bool {
	return strings.Contains(string(response), "\"users\":[]")
}*/
