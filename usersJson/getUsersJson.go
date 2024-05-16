package usersJson

//package bugz

/*import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	//. "phabricator/bugz"
)

func DownloadBugzillaUsers() {
	// Create a 'data/users' folder if it doesn't exist
	err := os.MkdirAll("data/users", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating 'users' folder: %v\n", err)
		return
	}

	// Extract unique IDs from bugs
	idUserMap, err := getIDS()
	if err != nil {
		return
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
}*/
