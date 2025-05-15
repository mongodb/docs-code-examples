package manage_indexes

type IndexExpectation struct {
	Name   string `bson:"name"`
	Fields []struct {
		Type          string `bson:"type"`
		Path          string `bson:"path"`
		NumDimensions int    `bson:"numDimensions"`
		Similarity    string `bson:"similarity"`
	}
}
