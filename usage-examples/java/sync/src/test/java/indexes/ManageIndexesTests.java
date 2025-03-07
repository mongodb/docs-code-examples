package indexes;

import com.fasterxml.jackson.databind.ObjectMapper;
import indexes.models.Field;
import indexes.models.IndexDefinition;
import org.bson.Document;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import java.util.ArrayList;
import static org.junit.jupiter.api.Assertions.*;

class ManageIndexesTests {
    @AfterEach
    void tearDown() {
        new DropIndex().main(new String[]{"Example placeholder arg"});
    }

    @Test
    @DisplayName("Test creating a basic Vector Search index")
    void TestCreateIndexBasic() {
        new CreateIndexBasic().main(new String[]{"Example placeholder arg"});
        Document index = ViewIndex.main(new String[]{"Example placeholder arg"});
        ObjectMapper objectMapper = new ObjectMapper();
        try {
            IndexDefinition indexAsObject = objectMapper.readValue(index.toJson(), IndexDefinition.class);
            ArrayList<Field> indexFields = indexAsObject.getLatestDefinition().getFields();
            Field definitionField = indexFields.get(0);
            assertEquals("vector", definitionField.getType());
            assertEquals("plot_embedding", definitionField.getPath());
            assertEquals(1536, definitionField.getNumDimensions());
            assertEquals("euclidean", definitionField.getSimilarity());
        } catch (Exception e) {
            fail("There was an error deserializing the index to an IndexDefinition " + e.getMessage());
        }
    }

    @Test
    @DisplayName("Test creating a Vector Search index with filter")
    void TestCreateIndexFilter() {
        new CreateIndexFilter().main(new String[]{"Example placeholder arg"});
        Document index = ViewIndex.main(new String[]{"Example placeholder arg"});
        ObjectMapper objectMapper = new ObjectMapper();
        try {
            IndexDefinition indexAsObject = objectMapper.readValue(index.toJson(), IndexDefinition.class);
            ArrayList<Field> indexFields = indexAsObject.getLatestDefinition().getFields();
            Field definitionField = indexFields.get(0);
            assertEquals("vector", definitionField.getType());
            assertEquals("plot_embedding", definitionField.getPath());
            assertEquals(1536, definitionField.getNumDimensions());
            assertEquals("euclidean", definitionField.getSimilarity());
            Field genresField = indexFields.get(1);
            assertEquals("filter", genresField.getType());
            assertEquals("genres", genresField.getPath());
            Field yearField = indexFields.get(2);
            assertEquals("filter", yearField.getType());
            assertEquals("year", yearField.getPath());
        } catch (Exception e) {
            fail("There was an error deserializing the index to an IndexDefinition " + e.getMessage());
        }
    }
}