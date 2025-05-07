#	:replace-start: {
#	  "terms": {
#	    "ATLAS_CONNECTION_STRING": "\"<connection-string>\""
#	  }
#	}
# 
# :snippet-start: example
from pymongo.mongo_client import MongoClient
import time # :remove:
import os # :remove:

def example():
    # Connect to your Atlas deployment
    ATLAS_CONNECTION_STRING = os.getenv("ATLAS_CONNECTION_STRING") # :remove:
    uri = ATLAS_CONNECTION_STRING
    client = MongoClient(uri)

    # Access your database and collection
    database = client["sample_mflix"]
    collection = database["embedded_movies"]

    index_name = "vector_index"
    # Delete your search index
    collection.drop_search_index(index_name)
    # :remove-start:
    """Wait to confirm the index is done deleting."""
    print("Polling to find out if the drop index operation is complete.")
    print("Note: this may take up to a minute.")
    predicate = None
    if predicate is None:
        predicate = lambda index: index.get("queryable") is True

    while True:
        indices = list(collection.list_search_indexes(index_name))
        if len(indices) and predicate(indices[0]):
            time.sleep(5)
        else:
            break

    print("Drop search index operation is complete.")
    # :remove-end:
    client.close()
# :snippet-end:
# :replace-end: