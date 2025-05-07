//	:replace-start: {
//	   "terms": {
//	      "manage_indexes": "main",
//	      "ExampleViewIndex(t *testing.T)": "main()"
//	   }
//	}
//
// :snippet-start: examples
package manage_indexes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv" // :remove:
	"log"
	"os"      // :remove:
	"testing" // :remove:

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExampleViewIndex(t *testing.T) {
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
	// Specify the options for the index to retrieve
	indexName := "vector_index"
	opts := options.SearchIndexes().SetName(indexName)
	// Get the index
	cursor, err := coll.SearchIndexes().List(ctx, opts)
	if err != nil {
		log.Fatalf("failed to get the index: %v", err)
	}
	// Print the index details to the console as JSON
	// :uncomment-start:
	// var results []bson.M
	// :uncomment-end:
	var results []IndexDefinition // :remove:
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatalf("failed to unmarshal results to bson: %v", err)
	}
	res, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("failed to marshal results to json: %v", err)
	}
	fmt.Println(string(res))
	// :remove-start:
	// The var below represents the commented-out 'results' above, which we don't actually use here for testing reasons
	var someBsonVar []bson.M
	fmt.Printf("Need to reference the variable and include BSON so the import shows in the examples%v\n", someBsonVar)
	if len(results) > 0 {
		fmt.Printf("Found %d indexes\n", len(results))
		fmt.Printf("This test should pass\n.")
	} else {
		t.Fail()
		fmt.Println("No indexes found.\n")
		fmt.Printf("This test should fail.\n")
	}
	// :remove-end:
}

// :snippet-end:
// :replace-end:
