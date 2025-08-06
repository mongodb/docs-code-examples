module atlas-sdk-go

go 1.24

// :remove-start:
// NOTE: confirm testify and indirect dependencies are removed in Bluehawk copy output
// once copied, confirm project builds successfully in artifact repo
// :remove-end:
require (
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0 // :remove:
	go.mongodb.org/atlas-sdk/v20250219001 v20250219001.1.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mongodb-forks/digest v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect; indirect // :remove:
	github.com/stretchr/objx v0.5.2 // indirect; indirect // :remove:
	golang.org/x/oauth2 v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
