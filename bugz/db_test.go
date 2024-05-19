package bugz

import (
	"testing"
	"zombiezen.com/go/sqlite/sqlitex"
)

func TestCreateAndInitializeDatabase(t *testing.T) {
	db, err := CreateAndInitializeDatabase("bugs.db")
	if err != nil {
		t.Fatalf("Failed to create and initialize database: %v", err)
	}
	defer db.Close()

	// Define the execOptions for the insert query
	execOptions := sqlitex.ExecOptions{
		Args: []interface{}{1, "2077-10-23 09:42:00", "John Dead", "Sample Bug", "{}"},
	}

	insertQuery, err := schemaFS.ReadFile("insert.sql")
	if err != nil {
		t.Fatalf("Failed to read insert query: %v", err)
	}

	if err := sqlitex.ExecuteTransient(db, string(insertQuery), &execOptions); err != nil {
		t.Fatalf("Error executing insert statement: %v", err)
	}

	// Read the select query from the embedded file
	selectQuery, err := schemaFS.ReadFile("select.sql")
	if err != nil {
		t.Fatalf("Failed to read select query: %v", err)
	}

	stmt, err := db.Prepare(string(selectQuery))

	if err != nil {
		t.Fatalf("Failed to prepare select statement: %v", err)
	}
	defer stmt.Finalize()

	// Execute the query and retrieve the results
	hasRow, err := stmt.Step()
	if err != nil {
		t.Fatalf("Error stepping through rows: %v", err)
	}
	if !hasRow {
		t.Fatalf("No rows found for id = 1")
	}

	id := stmt.ColumnInt64(0)
	creationTime := stmt.ColumnText(1)
	creator := stmt.ColumnText(2)
	summary := stmt.ColumnText(3)
	otherFieldsJSON := stmt.ColumnText(4)

	// Verify the retrieved data matches the sample data
	if id != 1 || creationTime != "2077-10-23 09:42:00" || creator != "John Dead" || summary != "Sample Bug" || otherFieldsJSON != "{}" {
		t.Fatalf("Retrieved data does not match sample data")
	}
}

func TestGetDistinctCreators(t *testing.T) {
	db, err := CreateAndInitializeDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create and initialize database: %v", err)
	}
	defer db.Close()

	sampleData := [][]interface{}{
		{1, "2077-10-23 09:42:00", "John Dead", "Sample Bug", "{}"},
		{2, "2077-10-23 10:00:00", "Jane Alive", "Another Bug", "{}"},
		{3, "2077-10-23 10:15:00", "John Dead", "Yet Another Bug", "{}"},
	}

	insertQuery, err := schemaFS.ReadFile("insert.sql")
	if err != nil {
		t.Fatalf("Failed to read insert query: %v", err)
	}

	for _, args := range sampleData {
		if err := sqlitex.ExecuteTransient(db, string(insertQuery), &sqlitex.ExecOptions{Args: args}); err != nil {
			t.Fatalf("Error executing insert statement: %v", err)
		}
	}

	creators, err := GetDistinctCreators(db)
	if err != nil {
		t.Fatalf("Failed to get distinct creators: %v", err)
	}

	expectedCreators := []string{"John Dead", "Jane Alive"}
	if len(creators) != len(expectedCreators) {
		t.Fatalf("Expected %v creators, got %v", len(expectedCreators), len(creators))
	}
	for i, creator := range creators {
		if creator != expectedCreators[i] {
			t.Fatalf("Expected creator %v, got %v", expectedCreators[i], creator)
		}
	}
}
