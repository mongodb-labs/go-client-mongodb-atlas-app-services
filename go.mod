module github.com/mongodb-labs/go-client-mongodb-atlas-app-services

go 1.23

toolchain go1.23.1

require (
	github.com/go-test/deep v1.1.1
	github.com/google/go-querystring v1.1.0
	go.mongodb.org/atlas v0.37.0
)

// Incorrect release versioning
retract v1.0.0
