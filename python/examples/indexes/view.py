#	:replace-start: {
#	  "terms": {
#	    "ATLAS_CONNECTION_STRING": "\"<connection-string>\""
#	  }
#	}
# 
# :snippet-start: example
from pymongo.mongo_client import MongoClient
import os # :remove:

def example():
    # Connect to your Atlas deployment
    ATLAS_CONNECTION_STRING = os.getenv("ATLAS_CONNECTION_STRING") # :remove:
    uri = ATLAS_CONNECTION_STRING
    client = MongoClient(uri)

    # Access your database and collection
    database = client["sample_mflix"]
    collection = database["embedded_movies"]

    # Get a list of the collection's search indexes and print them
    cursor = collection.list_search_indexes()
    docs = [] # :remove:
    for index in cursor:
        docs.append(index) # :remove:
        print(index)
    client.close()
    return docs # :remove:
# :snippet-end:
# :replace-end: