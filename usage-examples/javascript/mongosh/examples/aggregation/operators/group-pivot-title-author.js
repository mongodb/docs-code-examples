db.books.aggregate([
    { $group : { _id : "$author", books: { $push: "$title" } } }
])