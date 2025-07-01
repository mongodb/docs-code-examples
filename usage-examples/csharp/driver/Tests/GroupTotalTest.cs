using Examples.Aggregation;

namespace Tests;

public class GroupTotalTest
{
    private GroupTotal _example;
    [SetUp]
    public void Setup()
    {
        _example = new GroupTotal();
        _example.LoadSampleData();
    }

    [Test]
    public void TestOutputMatchesDocs()
    {
        var results = _example.PerformAggregation();
        
        var solutionRoot = DotNetEnv.Env.GetString("SOLUTION_ROOT", "Env variable not found. Verify you have a .env file with a valid connection string.");
        var outputLocation = "Examples/Aggregation/GroupTotalOutput.txt";
        var fullPath = Path.Combine(solutionRoot, outputLocation);
        var fileData = TestUtils.ReadBsonDocumentsFromFile(fullPath);
        
        Assert.That(results.Count, Is.EqualTo(fileData.Length), $"Result count {results.Count} does not match output example length {fileData.Length}.");
        for (var i = 0; i < fileData.Length; i++)  
        {  
            Assert.That(fileData[i], Is.EqualTo(results[i]), $"Mismatch at index {i}: expected {fileData[i]}, got {results[i]}.");
        }
    }
}