//	:replace-start: {
//	  "terms": {
//	    "process.env.ATLAS_CONNECTION_STRING": "\"<connection-string>\""
//	  }
//	}
// :snippet-start: example
import { MongoClient } from 'mongodb';

export async function createIndexFilter() {
    // connect to your Atlas deployment
    const uri =  process.env.ATLAS_CONNECTION_STRING;
    const client = new MongoClient(uri);
    try {
        const database = client.db("sample_mflix");
        const collection = database.collection("embedded_movies");

        // define your Atlas Vector Search index
        const indexName = "vector_index";
        const index = {
            name: indexName,
            type: "vectorSearch",
            definition: {
                "fields": [
                    {
                        "type": "vector",
                        "numDimensions": 1536,
                        "path": "plot_embedding",
                        "similarity": "euclidean"
                    },
                    {
                        "type": "filter",
                        "path": "genres"
                    },
                    {
                        "type": "filter",
                        "path": "year"
                    }
                ]
            }
        }

        // create the index
        const result = await collection.createSearchIndex(index);
        console.log(`Successfully created index named "${result}"`);

        // wait for the index to be ready to query
        console.log("Polling to confirm the index has finished building and can be queried.")
        console.log("NOTE: This may take up to a minute.")
        let isQueryable = false;
        while (!isQueryable) {
            const cursor = collection.listSearchIndexes();
            for await (const index of cursor) {
                if (index.name === indexName) {
                    if (index.queryable) {
                        console.log(`The search index "${indexName}" is queryable.`);
                        isQueryable = true;
                    } else {
                        await new Promise(resolve => setTimeout(resolve, 5000));
                    }
                }
            }
        }
    } finally {
        await client.close();
    }
}
// :uncomment-start:
//createIndexFilter().catch(console.dir);
// :uncomment-end:
// :snippet-end:
// :replace-end: