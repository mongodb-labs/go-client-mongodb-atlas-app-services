module github.com/mongodb-labs/go-client-mongodb-atlas-app-services

go 1.24

require (
	github.com/go-test/deep v1.1.2
	github.com/google/go-querystring v1.2.0
	go.mongodb.org/atlas v0.38.0
)

// Incorrect release versioning
retract v1.0.0
