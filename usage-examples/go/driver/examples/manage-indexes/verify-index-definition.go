package manage_indexes

import (
	"fmt"
	"strconv"
)

func VerifyIndexDefinition(results []IndexDefinition, expected []IndexExpectation) bool {
	localIsValid := true
	if len(results) != len(expected) {
		localIsValid = false
		return localIsValid
	}
	for i, result := range results {
		if len(result.LatestDefinition.Fields) != len(expected[i].Fields) {
			localIsValid = false
			fmt.Printf("Expected " + strconv.Itoa(len(expected[i].Fields)) + " fields in the index definition but got " + strconv.Itoa(len(result.LatestDefinition.Fields)) + " fields.\n")
			return localIsValid
		}
		if result.Name != expected[i].Name {
			localIsValid = false
			fmt.Printf("Expected the index name " + expected[i].Name + " but got " + result.Name + "\n")
		}

		for ii, expectedFields := range expected[i].Fields {
			if result.LatestDefinition.Fields[ii].Type != expectedFields.Type {
				localIsValid = false
				fmt.Printf("Expected the type " + expectedFields.Type + " but got " + result.LatestDefinition.Fields[ii].Type + "\n")
			}

			if result.LatestDefinition.Fields[ii].Path != expectedFields.Path {
				localIsValid = false
				fmt.Printf("Expected the path " + expectedFields.Path + " but got " + result.LatestDefinition.Fields[ii].Path + "\n")
			}

			if expectedFields.Type == "vectorSearch" {
				if result.LatestDefinition.Fields[ii].NumDimensions != expectedFields.NumDimensions {
					localIsValid = false
					fmt.Printf("Expected num dimensions to be %v, but got %v\n", strconv.Itoa(expectedFields.NumDimensions), strconv.Itoa(result.LatestDefinition.Fields[ii].NumDimensions))
				}

				if result.LatestDefinition.Fields[ii].Similarity != expectedFields.Similarity {
					localIsValid = false
					fmt.Printf("Expected the similarity " + expectedFields.Similarity + " but got " + result.LatestDefinition.Fields[ii].Similarity + "\n")
				}
			}
		}
	}
	return localIsValid
}
