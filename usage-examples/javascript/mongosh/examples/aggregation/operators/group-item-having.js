db.sales.aggregate(
    [
        // First Stage
        {
            $group :
                {
                    _id : "$item",
                    totalSaleAmount: { $sum: { $multiply: [ "$price", "$quantity" ] } }
                }
        },
        // Second Stage
        {
            $match: { "totalSaleAmount": { $gte: 100 } }
        }
    ]
)
