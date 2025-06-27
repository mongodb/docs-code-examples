const { exec, execSync } = require("child_process");
const fs = require("fs");
const path = require("path"); // To help with cross-platform file paths
const makeTempFileForTesting = require("../../utils/makeTempFileForTesting");
const unorderedArrayOutputMatches = require("../../utils/unorderedArrayOutputMatches");

jest.setTimeout(10000);

describe("mongosh group operator tests", () => {
    const mongoUri = process.env.CONNECTION_STRING;
    const port = process.env.CONNECTION_PORT;
    const dbName = "test";

    // Load test data before running the tests
    beforeAll(() => {
        const testDataSalesFilePath = "aggregation/operators/group-test-data-sales.js";
        const codeExampleDetailsSales = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: testDataSalesFilePath,
            validateOutput: true,
        }
        const tempTestDataSalesFilePath = makeTempFileForTesting(codeExampleDetailsSales);
        execSync(`mongosh --file ${tempTestDataSalesFilePath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }
        });

        const testDataBooksFilePath = "aggregation/operators/group-test-data-books.js";
        const codeExampleDetailsBooks = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: testDataBooksFilePath,
            validateOutput: true,
        }
        const tempTestDataBooksFilePath = makeTempFileForTesting(codeExampleDetailsBooks);
        execSync(`mongosh --file ${tempTestDataBooksFilePath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }
        });
    });

    // Drop the database and delete temp files after running the tests
    afterAll(() => {
        const command = `mongosh "${mongoUri}" --eval "db = db.getSiblingDB('${dbName}'); db.dropDatabase();"`

        try {
            execSync(command, { encoding: "utf8" })
        } catch (error) {
            console.error(`Failed to drop database '${dbName}':`, error.message);
        }

        const tempDirPath = path.resolve(__dirname, "../../temp");
        try {
            // Recursively delete the `/temp` directory and its contents
            fs.rmSync(tempDirPath, { recursive: true, force: true });
        } catch (error) {
            console.error("Failed to clean temp directory:", error.message);
        }
    });

    test("should group documents with a count", (done) => {
        const snippetFilePath = "aggregation/operators/group-count.js";

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Read the content of the expected output
        const expectedOutputFilePath = path.resolve(__dirname, "../../examples/aggregation/operators/group-count-output.json");
        const expectedOutput = fs.readFileSync(expectedOutputFilePath, "utf8");

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            expect(stdout).toContain(expectedOutput);
            done();
        });
    });

    test("should group documents with distinct values", (done) => {
        const snippetFilePath = "aggregation/operators/group-distinct.js";
        const expectedOutputFilePath = "aggregation/operators/group-distinct-output.json";

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            const result = unorderedArrayOutputMatches(expectedOutputFilePath, stdout)
            expect(result).toBe(true);
            done();
        });
    });

    test("should group documents by item field with condition", (done) => {
        const snippetFilePath = "aggregation/operators/group-item-having.js";
        const expectedOutputFilePath = "aggregation/operators/group-item-having-output.json";

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            const result = unorderedArrayOutputMatches(expectedOutputFilePath, stdout)
            expect(result).toBe(true);
            done();
        });
    });

    test("should group documents by day of year", (done) => {
        const snippetFilePath = "aggregation/operators/group-by-day.js";
        const expectedOutputFilePath = "aggregation/operators/group-by-day-output.json";

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            const result = unorderedArrayOutputMatches(expectedOutputFilePath, stdout)
            expect(result).toBe(true);
            done();
        });
    });

    test("should group documents by null", (done) => {
        const snippetFilePath = "aggregation/operators/group-by-null.js";
        const expectedOutputFilePath = path.resolve(__dirname, "../../examples/aggregation/operators/group-by-null-output.json");
        const expectedOutput = fs.readFileSync(expectedOutputFilePath, "utf8");

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            expect(stdout).toContain(expectedOutput);
            done();
        });
    });

    test("should group book titles by author", (done) => {
        const snippetFilePath = "aggregation/operators/group-pivot-title-author.js";
        const expectedOutputFilePath = "aggregation/operators/group-pivot-title-author-output.json";

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            const result = unorderedArrayOutputMatches(expectedOutputFilePath, stdout)
            expect(result).toBe(true);
            done();
        });
    });

    test("should group documents by author", (done) => {
        const snippetFilePath = "aggregation/operators/group-pivot-documents-author.js";
        const expectedOutputFilePath = "aggregation/operators/group-pivot-documents-author-output.json";

        const codeExampleDetails = {
            connectionString: mongoUri,
            dbName: dbName,
            filepath: snippetFilePath,
            validateOutput: true,
        }

        const tempScriptPath = makeTempFileForTesting(codeExampleDetails);

        // Execute the file using mongosh
        exec(`mongosh --file ${tempScriptPath} --port ${port}`, (error, stdout, stderr) => {
            expect(error).toBeNull(); // Ensure no error occurred
            if (stderr !== "") {
                console.error("Standard Error:", stderr);
            }

            // Validate the output
            const result = unorderedArrayOutputMatches(expectedOutputFilePath, stdout)
            expect(result).toBe(true);
            done();
        });
    });
});
