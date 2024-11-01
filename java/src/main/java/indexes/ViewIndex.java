//	:replace-start: {
//	  "terms": {
//      "Document": "void",
//	    "System.getenv(\"ATLAS_CONNECTION_STRING\")": "<connectionString>"
//	  }
//	}
package indexes;
// :snippet-start: example
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import org.bson.Document;

public class ViewIndex {
    public static Document main(String[] args) {
        // Replace the placeholder with your Atlas connection string
        String uri = System.getenv("ATLAS_CONNECTION_STRING");

        // Connect to your Atlas cluster
        try (MongoClient mongoClient = MongoClients.create(uri)) {

            // Set the namespace
            MongoDatabase database = mongoClient.getDatabase("sample_mflix");
            MongoCollection<Document> collection = database.getCollection("embedded_movies");

            // Specify the options for the index to retrieve
            String indexName = "vector_index";

            // Get the index and print details to the console as JSON
            try {
                Document listSearchIndex = collection.listSearchIndexes().name(indexName).first();
                if (listSearchIndex != null) {
                    System.out.println("Index found: " + listSearchIndex.toJson());
                    return listSearchIndex; // :remove:
                } else {
                    System.out.println("Index not found.");
                }
            } catch (Exception e) {
                throw new RuntimeException("Error finding index: " + e);
            }
        } catch (Exception e) {
            throw new RuntimeException("Error connecting to MongoDB: " + e);
        }
        return null; // :remove:
    }
}
// :snippet-end:
// :replace-end:
