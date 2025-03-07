package indexes.models;

public class Field {
    private String type;
    private String path;
    private Integer numDimensions;
    private String similarity;

    // Getters and Setters
    public String getType() {
        return type;
    }
    public void setType(String type) {
        this.type = type;
    }
    public String getPath() {
        return path;
    }
    public void setPath(String path) {
        this.path = path;
    }
    public Integer getNumDimensions() {
        return numDimensions;
    }
    public void setNumDimensions(Integer numDimensions) {
        this.numDimensions = numDimensions;
    }
    public String getSimilarity() {
        return similarity;
    }
    public void setSimilarity(String similarity) {
        this.similarity = similarity;
    }
}
