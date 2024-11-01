package tests

import (
	"os"
	"test-poc/examples/manage-indexes"
	"testing"
)

func TestCreateIndexBasic(t *testing.T) {
	// Test creating the index and performing a query that relies on the index
	manage_indexes.ExampleCreateIndexBasic(t)
	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

func TestCreateIndexFilter(t *testing.T) {
	// Test creating the index and performing a query that relies on the index
	manage_indexes.ExampleCreateIndexFilter(t)

	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

func TestViewIndex(t *testing.T) {
	manage_indexes.ExampleCreateIndexBasic(t)
	manage_indexes.ExampleViewIndex(t)

	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

func TestEditIndex(t *testing.T) {
	if os.Getenv("ENV") == "local" {
		t.Skip("Skipping this test in CI because it doesn't work on local Atlas")
	}
	manage_indexes.ExampleCreateIndexBasic(t)
	manage_indexes.ExampleEditIndex(t)

	// Drop the index to clear state for future tests
	manage_indexes.ExampleDropIndex()
}

func TestDropIndex(t *testing.T) {
	manage_indexes.ExampleCreateIndexBasic(t)
	manage_indexes.ExampleDropIndex()
}
