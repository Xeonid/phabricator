package bugz

import "testing"

func TestExtractIDs(t *testing.T) {
	// Define a Bug object with some data
	bug := Bug{
		AssignedToDetail: User{ID: 101},
		CCDetail:         []User{{ID: 102}, {ID: 103}},
		CreatorDetail:    User{ID: 104},
	}

	// Define the expected map of IDs
	expectedIDs := map[int]User{
		101: {ID: 101},
		102: {ID: 102},
		103: {ID: 103},
		104: {ID: 104},
	}

	// Extract IDs from the bug
	actualIDs := extractIDs(bug)

	// Check if the number of actual IDs matches the number of expected IDs
	if len(actualIDs) != len(expectedIDs) {
		t.Errorf("Number of actual IDs does not match number of expected IDs")
		return
	}

	// Compare individual elements
	for id, expectedUser := range expectedIDs {
		actualUser, ok := actualIDs[id]
		if !ok {
			t.Errorf("Expected ID %d not found in actual IDs", id)
			return
		}
		if actualUser.ID != expectedUser.ID {
			t.Errorf("Actual user ID %d does not match expected user ID %d for ID %d", actualUser.ID, expectedUser.ID, id)
		}
	}
}
