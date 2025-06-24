db.books.aggregate([
    // First Stage
    {
        $group : { _id : "$author", books: { $push: "$$ROOT" } }
    },
    // Second Stage
    {
        $addFields:
            {
                totalCopies : { $sum: "$books.copies" }
            }
    }
])