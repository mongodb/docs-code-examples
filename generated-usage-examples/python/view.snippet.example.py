from pymongo.mongo_client import MongoClient

def example():
    # Connect to your Atlas deployment
    uri = "<connection-string>"
    client = MongoClient(uri)

    # Access your database and collection
    database = client["sample_mflix"]
    collection = database["embedded_movies"]

    # Get a list of the collection's search indexes and print them
    cursor = collection.list_search_indexes()
    for index in cursor:
        print(index)
    client.close()
