const path = require("path");
const fs = require("fs");
const vm = require("vm");
const { Decimal128, ObjectId } = require("mongodb"); // Import MongoDB data types

function unorderedArrayOutputMatches(filepath, output) {
    // Read the content of the expected output
    const filepathString = "../examples/" + filepath;
    const outputFilePath = path.resolve(__dirname, filepathString);
    const expectedOutput = fs.readFileSync(outputFilePath, "utf8");

    // Define MongoDB types in the evaluation context and ensure proper handling
    const context = {
        Decimal128: (value) => new Decimal128(value), // Wrap Decimal128 class to include `new`
        ObjectId: (value) => new ObjectId(value), // Similarly wrap ObjectId
    };

    // Convert the expected output and actual output to arrays
    let expectedOutputArray, actualOutputArray;
    try {
        expectedOutputArray = vm.runInNewContext(expectedOutput, context); // Safely parse expected output
    } catch (error) {
        console.error("Failed to parse expected output:", error);
        return false;
    }

    try {
        actualOutputArray = vm.runInNewContext(output, context); // Safely parse actual output
    } catch (error) {
        console.error("Failed to parse actual output:", error);
        return false;
    }

    // Helper function to normalize MongoDB data types for comparison
    const normalizeItem = (item) => {
        const normalized = {};
        for (const key in item) {
            if (item[key] instanceof Decimal128 || item[key] instanceof ObjectId) {
                normalized[key] = item[key].toString(); // Convert Decimal128 and ObjectId to strings
            } else {
                normalized[key] = item[key]; // Keep other values as-is
            }
        }
        return normalized;
    };

    if (actualOutputArray !== undefined && expectedOutputArray !== undefined) {
        // Check that both arrays contain the same elements, regardless of order
        const isEqual = actualOutputArray.length === expectedOutputArray.length &&
            expectedOutputArray.every(expectedItem =>
                actualOutputArray.some(actualItem =>
                    JSON.stringify(normalizeItem(actualItem)) === JSON.stringify(normalizeItem(expectedItem))
                )
            );
        if (!isEqual) {
            console.log("Mismatch between actual output and expected output:", { actualOutputArray, expectedOutputArray });
        }

        return isEqual;
    } else {
        console.log("One or both arrays is undefined.");
        return false
    }
}

module.exports = unorderedArrayOutputMatches;