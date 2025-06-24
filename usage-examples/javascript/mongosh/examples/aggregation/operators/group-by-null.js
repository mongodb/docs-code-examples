db.sales.aggregate([
    {
        $group : {
            _id : null,
            totalSaleAmount: { $sum: { $multiply: [ "$price", "$quantity" ] } },
            averageQuantity: { $avg: "$quantity" },
            count: { $sum: 1 }
        }
    }
])