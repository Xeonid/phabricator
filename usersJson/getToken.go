package usersJson

//package bugz

/*import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetToken() (string, error) {

	apiURL := "https://bugs.freebsd.org/bugzilla/rest/login"

	// Create request body with login and password
	params := url.Values{}
	params.Set("login", "kovalthechosen@gmail.com")
	params.Set("password", "@tt.:3JmNeNh>7J")

	fullURL := apiURL + "?" + params.Encode()

	// Make a GET request to the API
	response, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making GET request to %s: %v\n", fullURL, err)
		return "", nil
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body from %s: %v\n", fullURL, err)
		return "", nil
	}
	response.Body.Close()

	// Define a struct to match the JSON structure
	var data map[string]interface{}

	// Unmarshal the JSON string into the struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", nil
	}

	// Retrieve the value associated with the key "token"
	tokenValue, ok := data["token"].(string)
	if !ok {
		fmt.Println("Token not found or not a string.")
		return "", nil
	}

	fmt.Println("Token Value:", tokenValue)
	return tokenValue, nil
}*/
