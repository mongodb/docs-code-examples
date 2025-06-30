namespace Examples.Aggregation;

using MongoDB.Driver;
using MongoDB.Bson;

public class GroupTotal
{
    private IMongoDatabase _aggDB;
    private IMongoCollection<Order> _orders;
    private string _appDir;
    private string _envRelPath;
    
    public void SeedData()
    {
        DotNetEnv.Env.TraversePath().Load();
        string uri = DotNetEnv.Env.GetString("CONNECTION_STRING", "Env variable not found. Verify you have a .env file with a valid connection string.");
        //var uri = "mongodb://localhost:27017";
        var client = new MongoClient(uri);
        _aggDB = client.GetDatabase("agg_tutorials_db");
        _orders = _aggDB.GetCollection<Order>("orders");
        _orders.DeleteMany(Builders<Order>.Filter.Empty);

        _orders.InsertMany(new List<Order>
        {
            new Order
            {
                CustomerId = "elise_smith@myemail.com",
                OrderDate = DateTime.Parse("2020-05-30T08:35:52Z"),
                Value = 231
            },
            new Order
            {
                CustomerId = "elise_smith@myemail.com",
                OrderDate = DateTime.Parse("2020-01-13T09:32:07Z"),
                Value = 99
            },
            new Order
            {
                CustomerId = "oranieri@warmmail.com",
                OrderDate = DateTime.Parse("2020-01-01T08:25:37Z"),
                Value = 63
            },
            new Order
            {
                CustomerId = "tj@wheresmyemail.com",
                OrderDate = DateTime.Parse("2019-05-28T19:13:32Z"),
                Value = 2
            },
            new Order
            {
                CustomerId = "tj@wheresmyemail.com",
                OrderDate = DateTime.Parse("2020-11-23T22:56:53Z"),
                Value = 187
            },
            new Order
            {
                CustomerId = "tj@wheresmyemail.com",
                OrderDate = DateTime.Parse("2020-08-18T23:04:48Z"),
                Value = 4
            },
            new Order
            {
                CustomerId = "elise_smith@myemail.com",
                OrderDate = DateTime.Parse("2020-12-26T08:55:46Z"),
                Value = 4
            },
            new Order
            {
                CustomerId = "tj@wheresmyemail.com",
                OrderDate = DateTime.Parse("2021-02-28T07:49:32Z"),
                Value = 1024
            },
            new Order
            {
                CustomerId = "elise_smith@myemail.com",
                OrderDate = DateTime.Parse("2020-10-03T13:49:44Z"),
                Value = 102
            }
        });
    }

    public List<BsonDocument> PerformAggregation()
    {
        DotNetEnv.Env.TraversePath().Load();
        string uri = DotNetEnv.Env.GetString("CONNECTION_STRING", "Env variable not found. Verify you have a .env file with a valid connection string.");
        //var uri = "mongodb://localhost:27017";
        var client = new MongoClient(uri);
        _aggDB = client.GetDatabase("agg_tutorials_db");
        _orders = _aggDB.GetCollection<Order>("orders");
        
        var results = _orders.Aggregate()
            .Match(o => o.OrderDate >= DateTime.Parse("2020-01-01T00:00:00Z") && o.OrderDate < DateTime.Parse("2021-01-01T00:00:00Z"))
            .SortBy(o => o.OrderDate)
            .Group(
                o => o.CustomerId,
                g => new
                {
                    CustomerId = g.Key,
                    FirstPurchaseDate = g.First().OrderDate,
                    TotalValue = g.Sum(i => i.Value),
                    TotalOrders = g.Count(),
                    Orders = g.Select(i => new { i.OrderDate, i.Value }).ToList()
                }
            )
            .SortBy(c => c.FirstPurchaseDate)
            .As<BsonDocument>();
            
        foreach (var result in results.ToList())
        {
            Console.WriteLine(result);
        }

        return results.ToList();
    }
}