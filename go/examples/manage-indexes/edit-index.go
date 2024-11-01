//	:replace-start: {
//	   "terms": {
//	      "manage_indexes": "main",
//	      "ExampleEditIndex(t *testing.T)": "main()"
//	   }
//	}
//
// :snippet-start: examples
package manage_indexes

import (
	"context"
	"fmt"
	"github.com/joho/godotenv" // :remove:
	"log"
	"os"      // :remove:
	"testing" // :remove:

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExampleEditIndex(t *testing.T) {
	ctx := context.Background()
	// :remove-start:
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}
	// Connect to your Atlas cluster
	uri := os.Getenv("ATLAS_CONNECTION_STRING")
	if uri == "" {
		log.Fatal("set your 'ATLAS_CONNECTION_STRING' environment variable.")
	}
	// :remove-end:
	// Replace the placeholder with your Atlas connection string
	// :uncomment-start:
	//const uri = "<connection-string>"
	// :uncomment-end:

	// Connect to your Atlas cluster
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to the server: %v", err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Set the namespace
	coll := client.Database("sample_mflix").Collection("embedded_movies")
	indexName := "vector_index"
	// :remove-start:
	// Get the index definition before performing the edit
	opts := options.SearchIndexes().SetName(indexName)
	cursor, err := coll.SearchIndexes().List(ctx, opts)
	if err != nil {
		log.Fatalf("failed to get the index: %v", err)
	}
	var results []IndexDefinition
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatalf("failed to unmarshal results to bson: %v", err)
	}
	numDimensionsBeforeEdit := results[0].LatestDefinition.Fields[0].NumDimensions
	// :remove-end:
	type vectorDefinitionField struct {
		Type          string `bson:"type"`
		Path          string `bson:"path"`
		NumDimensions int    `bson:"numDimensions"`
		Similarity    string `bson:"similarity"`
	}

	type vectorDefinition struct {
		Fields []vectorDefinitionField `bson:"fields"`
	}

	definition := vectorDefinition{
		Fields: []vectorDefinitionField{{
			Type:          "vector",
			Path:          "plot_embedding",
			NumDimensions: 1024,
			Similarity:    "euclidean"}},
	}
	err = coll.SearchIndexes().UpdateOne(ctx, indexName, definition)

	if err != nil {
		log.Fatalf("failed to update the index: %v", err)
	}

	fmt.Println("Successfully updated the search index")
	// :remove-start:
	// Get the index definition after performing the edit, and verify the change
	afterEditCursor, err := coll.SearchIndexes().List(ctx, opts)
	if err != nil {
		log.Fatalf("failed to get the index: %v", err)
	}
	var afterEditResults []IndexDefinition
	if err := afterEditCursor.All(ctx, &afterEditResults); err != nil {
		log.Fatalf("failed to unmarshal results to bson: %v", err)
	}
	numDimensionsAfterEdit := afterEditResults[0].LatestDefinition.Fields[0].NumDimensions
	if numDimensionsBeforeEdit != 1536 || numDimensionsAfterEdit != 1024 {
		t.Fail()
		fmt.Printf("the number of dimensions before the index edit is %d but expected 1536\n", numDimensionsBeforeEdit)
		fmt.Printf("the number of dimensions after the index edit is %d but expected 1024\n", numDimensionsAfterEdit)
		fmt.Printf("This test should fail.\n")
	} else {
		fmt.Printf("The number of dimensions before and after editing the index match the expected values.\n")
		fmt.Printf("This test should pass.\n")
	}
	// :remove-end:
}

// :snippet-end:
// :replace-end:
