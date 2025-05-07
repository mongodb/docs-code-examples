package indexes.models;

public class MainIndex {
    private String status;
    private Boolean queryable;
    private DefinitionVersion definitionVersion;
    private Definition definition;

    // Getters and Setters
    public String getStatus() {
        return status;
    }
    public void setStatus(String status) {
        this.status = status;
    }
    public Boolean getQueryable() {
        return queryable;
    }
    public void setQueryable(Boolean queryable) {
        this.queryable = queryable;
    }
    public DefinitionVersion getDefinitionVersion() {
        return definitionVersion;
    }
    public void setDefinitionVersion(DefinitionVersion definitionVersion) {
        this.definitionVersion = definitionVersion;
    }
    public Definition getDefinition() {
        return definition;
    }
    public void setDefinition(Definition definition) {
        this.definition = definition;
    }
}
