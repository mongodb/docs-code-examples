from pymongo.mongo_client import MongoClient
from pymongo.operations import SearchIndexModel
import time

def example():
  # Connect to your Atlas deployment
  uri = "<connection-string>"
  client = MongoClient(uri)

  # Access your database and collection
  database = client["sample_mflix"]
  collection = database["embedded_movies"]

  name="vector_index"

  # Create your index model, then create the search index
  search_index_model = SearchIndexModel(
    definition={
      "fields": [
        {
          "type": "vector",
          "path": "plot_embedding",
          "numDimensions": 1536,
          "similarity": "euclidean"
        }
      ]
    },
    name=name,
    type="vectorSearch",
  )

  result = collection.create_search_index(model=search_index_model)
  print("New search index named " + result + " is building.")

  """Wait for a search index to be ready."""
  print("Polling to find out if the search index is ready to query.")
  print("Note: this may take up to a minute.")
  predicate=None
  if predicate is None:
    predicate = lambda index: index.get("queryable") is True

  while True:
    indices = list(collection.list_search_indexes(name))
    if len(indices) and predicate(indices[0]):
      break
    time.sleep(5)
  print("Search index is ready to query.")

  client.close()
