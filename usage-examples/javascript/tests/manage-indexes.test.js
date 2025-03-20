import { createIndexBasic } from "../examples/indexes/create-index-basic.js";
import { createIndexFilter } from "../examples/indexes/create-index-filter.js";
import { viewIndex } from "../examples/indexes/view-index.js"
import { dropIndex } from "../examples/indexes/drop-index.js";

describe('Manage Indexes Tests', () => {
    it('Should return a definition for a basic vector index with no filter', async () => {
        await createIndexBasic();
        const indexes = await viewIndex();
        const vectorIndex = indexes[0];
        const latestDefinition = vectorIndex["latestDefinition"];
        const fields = latestDefinition.fields[0];
        expect(fields.type).toStrictEqual("vector");
        expect(fields.numDimensions).toStrictEqual(1536);
        expect(fields.path).toStrictEqual("plot_embedding");
        expect(fields.similarity).toStrictEqual("euclidean");
        await dropIndex();
    })
    it('Should return a definition for a vector index with filter', async () => {
        await createIndexFilter();
        const indexes = await viewIndex();
        const vectorIndex = indexes[0];
        const latestDefinition = vectorIndex["latestDefinition"];
        const vectorFields = latestDefinition.fields[0];
        expect(vectorFields.type).toStrictEqual("vector");
        expect(vectorFields.numDimensions).toStrictEqual(1536);
        expect(vectorFields.path).toStrictEqual("plot_embedding");
        expect(vectorFields.similarity).toStrictEqual("euclidean");
        const filterOne = latestDefinition.fields[1];
        expect(filterOne.type).toStrictEqual("filter");
        expect(filterOne.path).toStrictEqual("genres");
        const filterTwo = latestDefinition.fields[2];
        expect(filterTwo.type).toStrictEqual("filter");
        expect(filterTwo.path).toStrictEqual("year");
        await dropIndex();
    })
})
