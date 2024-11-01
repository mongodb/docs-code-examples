//	:replace-start: {
//	  "terms": {
//	    "System.getenv(\"ATLAS_CONNECTION_STRING\")": "<connectionString>"
//	  }
//	}
package indexes;
// :snippet-start: example
import com.mongodb.client.ListSearchIndexesIterable;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoCursor;
import com.mongodb.client.MongoDatabase;
import org.bson.Document;

public class DropIndex {
    public static void main(String[] args) {
        // Replace the placeholder with your Atlas connection string
        String uri = System.getenv("ATLAS_CONNECTION_STRING");

        // Connect to your Atlas cluster
        try (MongoClient mongoClient = MongoClients.create(uri)) {

            // Set the namespace
            MongoDatabase database = mongoClient.getDatabase("sample_mflix");
            MongoCollection<Document> collection = database.getCollection("embedded_movies");

            // Specify the index to delete
            String indexName = "vector_index";

            try {
                collection.dropSearchIndex(indexName);
            } catch (Exception e) {
                throw new RuntimeException("Error deleting index: " + e);
            }
            // :remove-start:
            // Wait for the drop index operation to complete
            System.out.println("Polling to confirm the index has successfully been deleted.");

            ListSearchIndexesIterable<Document> searchIndexes = collection.listSearchIndexes();
            boolean isDeleted = false;
            Document doc = null;
            while (!isDeleted) {
                boolean indexMatchingNameExists = false;
                ListSearchIndexesIterable<Document> innerSearchIndexes = collection.listSearchIndexes();
                try (MongoCursor<Document> cursor = innerSearchIndexes.iterator()) {
                    if (!cursor.hasNext()) {
                        break;
                    }
                    Document current = cursor.next();
                    String name = current.getString("name");
                    if (name.equals(indexName)) {
                        indexMatchingNameExists = true;
                    }
                    if (indexMatchingNameExists) {
                        Thread.sleep(500);
                    } else {
                        isDeleted = true;
                    }
                } catch (Exception e) {
                    throw new RuntimeException("Failed to list search indexes: " + e);
                }
            }
            System.out.println(indexName + " index has successfully been deleted.");
            // :remove-end:
        } catch (Exception e) {
            throw new RuntimeException("Error connecting to MongoDB: " + e);
        }
    }
}
// :snippet-end:
// :replace-end:
