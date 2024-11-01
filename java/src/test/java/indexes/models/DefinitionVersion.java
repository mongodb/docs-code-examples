package indexes.models;

public class DefinitionVersion {
    private Integer version;
    private MongoDate createdAt;

    // Getters and Setters
    public Integer getVersion() {
        return version;
    }
    public void setVersion(Integer version) {
        this.version = version;
    }
    public MongoDate getCreatedAt() {
        return createdAt;
    }
    public void setCreatedAt(MongoDate createdAt) {
        this.createdAt = createdAt;
    }
}
