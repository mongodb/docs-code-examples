package tests

import (
	"fmt"
	run_queries "test-poc/examples/run-queries"
)

func VerifyMovieQueryOutput(results []run_queries.ProjectedMovieResult, expected []run_queries.ProjectedMovieResult) bool {
	localIsValid := true
	if len(results) != len(expected) {
		localIsValid = false
		fmt.Printf("Expected %v results, got %v results.\n", len(expected), len(results))
		fmt.Printf("There's a mismatch between the number of results, so this test should fail.\n")
		return localIsValid
	}
	for i, result := range results {

		if result != expected[i] {
			if result.Title != expected[i].Title {
				fmt.Printf("Title: Got \"%v\" and expected \"%v\"\n", result.Title, expected[i].Title)
			}
			if result.Plot != expected[i].Plot {
				fmt.Printf("Plot: For %v, got \"%v\" and expected \"%v\"\n", result.Title, result.Plot, expected[i].Plot)
			}
			if result.Score != expected[i].Score {
				fmt.Printf("Score: For %v, got \"%v\" and expected \"%v\"\n", result.Title, result.Score, expected[i].Score)
			}
			localIsValid = false
		}
	}
	return localIsValid
}
