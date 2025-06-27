db.sales.aggregate([
    // First Stage
    {
        $match : { "date": { $gte: new ISODate("2014-01-01"), $lt: new ISODate("2015-01-01") } }
    },
    // Second Stage
    {
        $group : {
            _id : { $dateToString: { format: "%Y-%m-%d", date: "$date" } },
            totalSaleAmount: { $sum: { $multiply: [ "$price", "$quantity" ] } },
            averageQuantity: { $avg: "$quantity" },
            count: { $sum: 1 }
        }
    },
    // Third Stage
    {
        $sort : { totalSaleAmount: -1 }
    }
])