package main

import (
	. "phabricator/bugz"
	//. "phabricator/usersJson"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments. use help")
		return
	}
	command := os.Args[1]

	databasePath := "bugsNew.db"      // Specify the path to the database
	bc := NewBugzClient(databasePath) // Create a BugzClient instance

	switch command {
	case "bugzilla-download-bugs":
		err := bc.DownloadBugzillaBugs()
		if err != nil {
			fmt.Printf("Error downloading bugs: %v\n", err)
		}
	case "help":
		printHelp()
	case "bugzilla-list-bugs":
		err := bc.ListBugs()
		if err != nil {
			fmt.Printf("Error listing bugs: %v\n", err)
		}
	case "bugzilla-download-users":
		err := bc.DownloadBugzillaUsers(databasePath)
		if err != nil {
			fmt.Printf("Error downloading bugs: %v\n", err)
		}
	case "bugzilla-show-bugs":
		if err := bc.ShowBugs(); err != nil {
			fmt.Printf("error showing bugs: %v\n", err)
		}
	default:
		fmt.Println("invalid command")
		printHelp()
	}

}

func printHelp() { // not sure about functions descriptions
	fmt.Println("available commands:\n" +
		"bugzilla-download-bugs - downloads bugs from bugzilla\n" +
		"bugzilla-show-bugs - shows bugzilla bugs\n" +
		"bugzilla-list-bugs - displays downloaded bugs\n" +
		"bugzilla-download-users - downloads users from bugzilla\n" +
		"help - shows available commands")
}
