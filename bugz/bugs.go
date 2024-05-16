package bugz

/*import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func DownloadBugzillaBugs() error {
	apiURL := "https://bugs.freebsd.org/bugzilla/rest/bug"

	// Create a 'bugs' folder if it doesn't exist
	err := os.MkdirAll("bugs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error creating 'users' folder: %v\n", err)
	}

	// Specify the pagination parameters
	pageSize := 1000
	pageNumber := 0

	for {
		// Create query parameters
		params := url.Values{}
		params.Set("api_key", "UVy7lq2agakMQHafj1p2dKaTy2gcxM4t4P1cYvUH")
		params.Set("limit", fmt.Sprint(pageSize))
		params.Set("offset", fmt.Sprint((pageNumber)*pageSize))

		// Construct the full URL with query parameters
		fullURL := apiURL + "?" + params.Encode()

		// Make a GET request to the API
		response, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("Error making GET request to %s: %v\n", fullURL, err)
		}

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body from %s: %v\n", fullURL, err)
		}
		response.Body.Close()

		// Process the JSON data
		var bugsResponse map[string][]Bug
		err = json.Unmarshal(body, &bugsResponse)
		if err != nil {
			return fmt.Errorf("Error decoding JSON: %v\n", err)
		}

		// Iterate over bugs and write to individual files
		for _, bug := range bugsResponse["bugs"] {
			filename := filepath.Join("bugsJson", fmt.Sprintf("bug_%d.json", bug.ID))
			err := writeToFile(filename, bug)
			if err != nil {
				return fmt.Errorf("Error writing to file %s: %v\n", filename, err)
			}
		}

		// Check if there are more pages
		if len(bugsResponse["bugs"]) < pageSize {
			break
		}

		// Move to the next page
		pageNumber++
	}

	fmt.Println("Bug data written to individual files.")
	return nil
}

func writeToFile(filename string, bug Bug) error {
	// Use json.MarshalIndent to preserve the original formatting
	indentedData, err := json.MarshalIndent(bug, "", "  ")
	if err != nil {
		return err
	}

	// Write the data to the specified file
	return os.WriteFile(filename, indentedData, 0644)
}*/
