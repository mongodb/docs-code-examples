import { MongoClient } from 'mongodb';

export async function viewIndex() {
    // connect to your Atlas deployment
    const uri =  "<connection-string>";
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
viewIndex().catch(console.dir);
