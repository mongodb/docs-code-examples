package indexes.models;

public class StatusDetail {
    private String hostname;
    private String status;
    private Boolean queryable;
    private MainIndex mainIndex;

    // Getters and Setters
    public String getHostname() {
        return hostname;
    }
    public void setHostname(String hostname) {
        this.hostname = hostname;
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
    public MainIndex getMainIndex() {
        return mainIndex;
    }
    public void setMainIndex(MainIndex mainIndex) {
        this.mainIndex = mainIndex;
    }
}
