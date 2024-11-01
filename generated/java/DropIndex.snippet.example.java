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
        String uri = <connectionString>;

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
        } catch (Exception e) {
            throw new RuntimeException("Error connecting to MongoDB: " + e);
        }
    }
}
