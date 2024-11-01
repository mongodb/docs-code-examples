//	:replace-start: {
//	  "terms": {
//	    "manage_indexes": "main",
//	    "ExampleCreateIndexBasic(t *testing.T)": "main()"
//	  }
//	}
//
// :snippet-start: examples
package manage_indexes

import (
	"context"
	"fmt"
	"log"
	"os"      // :remove:
	"testing" // :remove:
	"time"

	"github.com/joho/godotenv" // :remove:

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExampleCreateIndexBasic(t *testing.T) {
	ctx := context.Background()
	// :remove-start:
	// Get the connection string from the .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}
	uri := os.Getenv("ATLAS_CONNECTION_STRING")
	if uri == "" {
		log.Fatal("set your 'ATLAS_CONNECTION_STRING' environment variable.")
	}
	// :remove-end:
	// Replace the placeholder
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

	// Define the index details
	type vectorDefinitionField struct {
		Type          string `bson:"type"`
		Path          string `bson:"path"`
		NumDimensions int    `bson:"numDimensions"`
		Similarity    string `bson:"similarity"`
	}

	type vectorDefinition struct {
		Fields []vectorDefinitionField `bson:"fields"`
	}

	indexName := "vector_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")

	indexModel := mongo.SearchIndexModel{
		Definition: vectorDefinition{
			Fields: []vectorDefinitionField{{
				Type:          "vector",
				Path:          "plot_embedding",
				NumDimensions: 1536,
				Similarity:    "euclidean"}},
		},
		Options: opts,
	}

	// Create the index
	fmt.Println("Creating the index.")
	searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("failed to create the search index: %v", err)
	}

	// Await the creation of the index.
	fmt.Println("Polling to confirm successful index creation.")
	fmt.Println("NOTE: This may take up to a minute.")
	searchIndexes := coll.SearchIndexes()
	var doc bson.Raw
	for doc == nil {
		cursor, err := searchIndexes.List(ctx, options.SearchIndexes().SetName(searchIndexName))
		if err != nil {
			fmt.Errorf("failed to list search indexes: %w", err)
		}

		if !cursor.Next(ctx) {
			break
		}

		name := cursor.Current.Lookup("name").StringValue()
		queryable := cursor.Current.Lookup("queryable").Boolean()
		if name == searchIndexName && queryable {
			doc = cursor.Current
			// :remove-start:
			var definitions []IndexDefinition
			if err := cursor.All(ctx, &definitions); err != nil {
				log.Fatalf("failed to unmarshal results to IndexDefinitions: %v", err)
			}
			expected := IndexExpectation{
				Name: "vector_index",
				Fields: []struct {
					Type          string `bson:"type"`
					Path          string `bson:"path"`
					NumDimensions int    `bson:"numDimensions"`
					Similarity    string `bson:"similarity"`
				}{{"vector", "plot_embedding", 1536, "euclidean"}},
			}
			if VerifyIndexDefinition(definitions, []IndexExpectation{expected}) {
				fmt.Printf("The relevant parts of the index definition match the expected outputs.\n")
				fmt.Printf("This test should pass.\n")
			} else {
				t.Fail()
				fmt.Printf("The relevant parts of the index definition do not match the expected outputs.\n")
				fmt.Printf("This test should fail.\n")
			}
			// :remove-end:
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Println("Name of Index Created: " + searchIndexName)
}

// :snippet-end:
// :replace-end:
