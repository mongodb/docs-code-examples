package indexes.models;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Date;

public class MongoDate {
    @JsonProperty("$date")
    Date date;
}