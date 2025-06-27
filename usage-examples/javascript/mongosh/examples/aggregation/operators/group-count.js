db.sales.aggregate( [
    {
        $group: {
            _id: null,
            count: { $count: { } }
        }
    }
] )