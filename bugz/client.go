package bugz

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type BugzClient struct {
	token string
	URL   string

	http *http.Client
}

type BugzLoginResponse struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

func NewBugzClient() *BugzClient {
	login := os.Getenv("BUGZILLA_LOGIN") //Retrieve env var values and check if they are empty
	password := os.Getenv("BUGZILLA_PASSWORD")
	if login == "" || password == "" {
		panic("BUGZILLA_LOGIN or BUGZILLA_PASSWORD is not set")
	}

	bc := &BugzClient{
		URL:   "https://bugs.freebsd.org/bugzilla/rest",
		token: "",
		http:  &http.Client{},
	}

	formData := url.Values{}
	formData.Set("login", login)
	formData.Set("password", password)

	response, err := bc.http.Get(bc.URL + "/login?" + formData.Encode())
	if err != nil {
		fmt.Printf("login and/or password incorrect")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("login failed, status code: %d", response.StatusCode)
	}

	var loginResponse BugzLoginResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		fmt.Printf("error reading bugzilla login response body: %s", err)
	}

	if loginResponse.Token == "" {
		fmt.Printf("login token is empty")
	}
	bc.token = loginResponse.Token

	return bc
}

// DownloadBugzillaBugs downloads all bugs from the Bugzilla API and saves them to individual JSON files.
func (bc *BugzClient) DownloadBugzillaBugs() error {
	// Make URL to bugs
	apiURL := bc.URL + "/bug"

	// Create a 'bugs' folder if it doesn't exist
	err := os.MkdirAll("bugs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating 'users' folder: %v", err)
	}

	// Specify the pagination parameters
	pageSize := 1000
	pageNumber := 0

	for {
		// Create query parameters
		params := url.Values{}
		params.Set("token", bc.token)
		params.Set("limit", fmt.Sprint(pageSize))
		params.Set("offset", fmt.Sprint((pageNumber)*pageSize))

		// Construct the full URL with query parameters
		fullURL := apiURL + "?" + params.Encode()

		// Make a GET request to the API
		response, err := bc.http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error making GET request to %s: %v", fullURL, err)
		}

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error reading response body from %s: %v", fullURL, err)
		}
		response.Body.Close()

		// Process the JSON data
		var bugsResponse map[string][]Bug
		err = json.Unmarshal(body, &bugsResponse)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}

		// Iterate over bugs and write to individual files
		for _, bug := range bugsResponse["bugs"] {
			filename := filepath.Join("bugs", fmt.Sprintf("bug_%d.json", bug.ID))
			err := writeToFile(filename, bug)
			if err != nil {
				return fmt.Errorf("error writing to file %s: %v", filename, err)
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
}

func (bc *BugzClient) ShowBugs() error {
	fmt.Println("showing bugs")
	return nil
}

func (bc *BugzClient) ListBugs() error {
	fmt.Println("listing bugs")
	return nil
}

func getIDS() (map[int]User, error) {
	// Specify the path to the directory containing JSON files
	directoryPath := "bugs"

	fmt.Println("Get IDS from bugs")

	// Get a list of all JSON files in the directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	// Map to store extracted IDs and corresponding User objects
	idUserMap := make(map[int]User)

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

			// Add the extracted IDs to the map
			ids := extractIDs(bug)
			for id, user := range ids {
				idUserMap[id] = user
			}
		}
	}

	return idUserMap, nil
}

func extractIDs(bug Bug) map[int]User {
	// Create a map to store the unique IDs and corresponding Users
	idUserMap := make(map[int]User)

	// Extract 'id' field from AssignedToDetail
	if _, ok := idUserMap[bug.AssignedToDetail.ID]; !ok {
		idUserMap[bug.AssignedToDetail.ID] = bug.AssignedToDetail
	}

	// Extract 'id' field from CCDetail
	for _, ccUser := range bug.CCDetail {
		if _, ok := idUserMap[ccUser.ID]; !ok {
			idUserMap[ccUser.ID] = ccUser
		}
	}

	// Extract 'id' field from CreatorDetail
	if _, ok := idUserMap[bug.CreatorDetail.ID]; !ok {
		idUserMap[bug.CreatorDetail.ID] = bug.CreatorDetail
	}

	return idUserMap
}

func (bc *BugzClient) DownloadBugzillaUsers() error {
	// Create a 'data/users' folder if it doesn't exist
	err := os.MkdirAll("data/users", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating 'users' folder: %w", err)
	}

	// Extract unique IDs from bugs
	idUserMap, err := getIDS()
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	// Iterate through the map
	for id, user := range idUserMap {
		// Construct the filename
		filename := fmt.Sprintf("user_%d.json", id)

		// Create the full file path
		filePath := filepath.Join("data/users", filename)

		// Open the file for writing
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			continue
		}
		defer file.Close()

		// Encode the user object to JSON and write it to the file
		encoder := json.NewEncoder(file)
		err = encoder.Encode(user)
		if err != nil {
			fmt.Printf("Error encoding user %d to JSON: %v\n", id, err)
			continue
		}

		fmt.Printf("User %d saved to %s\n", id, filePath)
	}
	return nil
}

//go:embed schema.sql
var schemaFS embed.FS

//go:embed queries.sql
var queriesFS embed.FS

func CreateAndInitializeDatabase(databasePath string) (*sqlite.Conn, error) {
	db, err := sqlite.OpenConn(databasePath, 0)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Read the schema from the embedded file
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema: %v", err)
	}

	if err := sqlitex.ExecScript(db, string(schema)); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	return db, nil
}

func GetDistinctCreators(db *sqlite.Conn) ([]string, error) {
	query, err := queriesFS.ReadFile("queries.sql")
	if err != nil {
		log.Fatalf("Failed to read query: %v", err)
	}

	stmt, err := db.Prepare(string(query))
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Finalize()

	var creators []string
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			log.Fatalf("Error stepping through rows: %v", err)
		}
		if !hasRow {
			break
		}
		creators = append(creators, stmt.ColumnText(0))
	}

	return creators, nil
}
