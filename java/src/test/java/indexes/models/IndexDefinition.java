package indexes.models;

import org.bson.types.ObjectId;

import java.util.ArrayList;

public class IndexDefinition {
    private ObjectId id;
    private String name;
    private String type;
    private String status;
    private Boolean queryable;
    private Integer latestVersion;
    private DefinitionVersion latestDefinitionVersion;
    private Definition latestDefinition;
    private ArrayList<StatusDetail> statusDetail;

    // Getters and Setters
    public ObjectId getId() {
        return id;
    }
    public void setId(ObjectId id) {
        this.id = id;
    }
    public String getName() {
        return name;
    }
    public void setName(String name) {
        this.name = name;
    }
    public String getType() {
        return type;
    }
    public void setType(String type) {
        this.type = type;
    }
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
    public Integer getLatestVersion() { return latestVersion; };
    public void setLatestVersion(Integer latestVersion) { this.latestVersion = latestVersion; };
    public DefinitionVersion getLatestDefinitionVersion() {
        return latestDefinitionVersion;
    }
    public void setLatestDefinitionVersion(DefinitionVersion latestDefinitionVersion) {
        this.latestDefinitionVersion = latestDefinitionVersion;
    }
    public Definition getLatestDefinition() {
        return latestDefinition;
    }
    public void setLatestDefinition(Definition latestDefinition) {
        this.latestDefinition = latestDefinition;
    }
    public ArrayList<StatusDetail> getStatusDetail() {
        return statusDetail;
    }
    public void setStatusDetail(ArrayList<StatusDetail> statusDetail) {
        this.statusDetail = statusDetail;
    }
}

