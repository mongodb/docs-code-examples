from pymongo.mongo_client import MongoClient

def example():
    # Connect to your Atlas deployment
    uri = "<connection-string>"
    client = MongoClient(uri)

    # Access your database and collection
    database = client["sample_mflix"]
    collection = database["embedded_movies"]

    index_name = "vector_index"
    # Delete your search index
    collection.drop_search_index(index_name)
    client.close()
