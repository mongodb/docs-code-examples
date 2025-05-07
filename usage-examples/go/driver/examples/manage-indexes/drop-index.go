//	:replace-start: {
//	   "terms": {
//	      "manage_indexes": "main",
//	      "ExampleDropIndex": "main"
//	   }
//	}
//
// :snippet-start: examples
package manage_indexes

import (
	"context"
	"fmt"
	"log"
	"os"   // :remove:
	"time" // :remove:

	"github.com/joho/godotenv" // :remove:

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExampleDropIndex() {
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

	err = coll.SearchIndexes().DropOne(ctx, indexName)
	if err != nil {
		log.Fatalf("failed to delete the index: %v", err)
	}

	fmt.Println("Successfully deleted the Vector Search index")
	// :remove-start:
	fmt.Println("Polling to confirm successful index deletion.")
	fmt.Println("NOTE: This may take up to a minute.")
	searchIndexes := coll.SearchIndexes()
	indexNotYetDeleted := true
	loopNumber := 0
	for indexNotYetDeleted {
		cursor, err := searchIndexes.List(ctx, options.SearchIndexes().SetName(indexName))
		if err != nil {
			fmt.Errorf("failed to list search indexes: %w", err)
		}

		if !cursor.Next(ctx) {
			break
		}

		name := cursor.Current.Lookup("name").StringValue()
		// If dropping the index takes more than a minute, which is 12 loops
		// with a 5 second sleep, something has gone wrong and we should
		// abandon all hope
		if name == indexName && loopNumber < 12 {
			time.Sleep(5 * time.Second)
			loopNumber += 1
		} else if name == indexName && loopNumber == 12 {
			log.Fatalf("Attempted to drop the index for a minute but it's still there. Something went wrong.")
			indexNotYetDeleted = false
		} else {
			indexNotYetDeleted = false
		}
	}
	fmt.Println("Index named " + indexName + " was deleted.")
	// :remove-end:
}

// :snippet-end:
// :replace-end:
