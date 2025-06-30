using Examples.Aggregation;

namespace Tests;

public class GroupTotalTest
{
    [SetUp]
    public void Setup()
    {
        var obj = new GroupTotal();
        obj.SeedData();
    }

    [Test]
    public void TestOutputMatchesDocs()
    {
        var obj = new GroupTotal();
        var results = obj.PerformAggregation();
        
        DotNetEnv.Env.TraversePath().Load();
        string solutionRoot = DotNetEnv.Env.GetString("SOLUTION_ROOT", "Env variable not found. Verify you have a .env file with a valid connection string.");
        string outputLocation = "Examples/Aggregation/GroupTotalOutput.txt";
        string fullPath = Path.Combine(solutionRoot, outputLocation);
        var fileData = TestUtils.ReadBsonDocumentsFromFile(fullPath);
        
        Assert.That(results.Count, Is.EqualTo(fileData.Length), "Result count does not match output example length.");
        for (int i = 0; i < fileData.Length; i++)  
        {  
            Assert.That(fileData[i], Is.EqualTo(results[i]), $"Mismatch at index {i}: expected {fileData[i]}, got {results[i]}.");
        }
    }
}