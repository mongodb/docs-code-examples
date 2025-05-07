import { MongoClient } from 'mongodb';

export async function dropIndex() {
    // connect to your Atlas deployment
    const uri =  "<connection-string>";
    const client = new MongoClient(uri);
    try {
        const database = client.db("sample_mflix");
        const collection = database.collection("embedded_movies");
        const indexName = "vector_index";
        await collection.dropSearchIndex(indexName);
        console.log(`Successfully dropped the index named "${indexName}"`);

        // wait for the index to be deleted from MongoDB
        console.log("Polling to confirm the index drop operation has completed.")
        console.log("NOTE: This may take up to a minute.")
        let isDeleted = false;
        while (!isDeleted) {
            const cursor = collection.listSearchIndexes();
            let indexMatchingNameExists = false;
            for await (const index of cursor) {
                if (index.name === indexName) {
                    indexMatchingNameExists = true;
                }
            }
            if (indexMatchingNameExists) {
                await new Promise(resolve => setTimeout(resolve, 5000));
            } else {
                console.log(`The search index "${indexName}" is deleted.`);
                isDeleted = true;
            }
        }
    } finally {
        await client.close();
    }
}
dropIndex().catch(console.dir);
