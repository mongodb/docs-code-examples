//	:replace-start: {
//	  "terms": {
//	    "process.env.ATLAS_CONNECTION_STRING": "\"<connection-string>\""
//	  }
//	}
// :snippet-start: example
import { MongoClient } from 'mongodb';

export async function viewIndex() {
    // connect to your Atlas deployment
    const uri =  process.env.ATLAS_CONNECTION_STRING;
    const client = new MongoClient(uri);
    try {
        const database = client.db("sample_mflix");
        const collection = database.collection("embedded_movies");

        // run the helper method
        const result = await collection.listSearchIndexes("vector_index");
        let indexes = [];
        for await (const index of result) {
            console.log(index);
            indexes.push(index);
        }
        return indexes;
    } finally {
        await client.close();
    }
}
// :uncomment-start:
//viewIndex().catch(console.dir);
// :uncomment-end:
// :snippet-end:
// :replace-end: