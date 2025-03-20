package manage_indexes

import "go.mongodb.org/mongo-driver/bson/primitive"

type IndexDefinition struct {
	Id                      string `bson:"id"`
	Name                    string `bson:"name"`
	Type                    string `bson:"type"`
	Status                  string `bson:"status"`
	Queryable               bool   `bson:"queryable"`
	LatestDefinitionVersion struct {
		Version   int64              `bson:"version"`
		CreatedAt primitive.DateTime `bson:"createdAt"`
	} `bson:"latestDefinitionVersion"`
	LatestDefinition struct {
		Fields []struct {
			Type          string `bson:"type"`
			Path          string `bson:"path"`
			NumDimensions int    `bson:"numDimensions"`
			Similarity    string `bson:"similarity"`
		} `bson:"fields"`
	} `bson:"latestDefinition"`
	StatusDetail []struct {
		Hostname  string `bson:"hostname"`
		Status    string `bson:"status"`
		Queryable bool   `bson:"queryable"`
		MainIndex struct {
			Status            string `bson:"status"`
			Queryable         bool   `bson:"queryable"`
			DefinitionVersion struct {
				Version   int64              `bson:"version"`
				CreatedAt primitive.DateTime `bson:"createdAt"`
			} `bson:"definitionVersion"`
			Definition struct {
				Fields []struct {
					Type          string `bson:"type"`
					Path          string `bson:"path"`
					NumDimensions int64  `bson:"numDimensions,omitempty"`
					Similarity    string `bson:"similarity,omitempty"`
				} `bson:"fields"`
			} `bson:"definition"`
		} `bson:"mainIndex"`
	} `bson:"statusDetail"`
}
