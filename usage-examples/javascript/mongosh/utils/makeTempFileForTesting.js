const fs = require("fs");
const path = require("path"); // To help with cross-platform file paths

function makeTempFileForTesting(details) {
    // Read content of code snippet
    const filepathString = "../examples/" + details.filepath;
    const snippetFilePath = path.resolve(__dirname, filepathString);
    const codeSnippet = fs.readFileSync(snippetFilePath, "utf8");

    // Construct code example string with details to test the code example snippet text
    // Add connection method
    let tempFileContents = `db = connect('${details.connectionString}`;

    // If there's a DB name to specify, add it to the connection string
    if (details.dbName) {
        tempFileContents = `${tempFileContents}/${details.dbName}');`;
    }

    // If we want to validate the output of the example, use the `printjson` method - otherwise, don't bother
    if (details.validateOutput) {
        tempFileContents = `${tempFileContents}\nprintjson(${codeSnippet});`;
    } else {
        tempFileContents = `${tempFileContents}\n${codeSnippet};`;
    }

    // Write code snippet + setup code to a temporary file
    const tempDir = "../temp";
    const buildTempFilepath = `${tempDir}/${details.filepath}`;
    const tempScriptPath = path.resolve(__dirname, buildTempFilepath);
    const tempScriptDir = path.dirname(tempScriptPath); // Extract the directory portion of the file path

    // Create the temp directory and its parent directories recursively if they don't exist
    fs.mkdirSync(tempScriptDir, { recursive: true });

    fs.writeFileSync(tempScriptPath, tempFileContents);
    return tempScriptPath;
}

module.exports = makeTempFileForTesting;