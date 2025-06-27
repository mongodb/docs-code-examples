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
    public void Test1()
    {
        var obj = new GroupTotal();
        var results = obj.PerformAggregation();
        
        string outputfilePath = "/Users/dachary.carey/workspace/docs-code-examples/usage-examples/csharp/driver/Examples/Aggregation/GroupTotalOutput.txt";
        var fileData = TestUtils.ReadBsonDocumentsFromFile(outputfilePath);
        Assert.That(results.Count, Is.EqualTo(fileData.Length), "Result count does not match output example length.");
        for (int i = 0; i < fileData.Length; i++)  
        {  
            Assert.That(fileData[i], Is.EqualTo(results[i]), $"Mismatch at index {i}: expected {fileData[i]}, got {results[i]}.");
        }
    }
}