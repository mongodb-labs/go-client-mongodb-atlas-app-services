# go-client-mongodb-atlas-app-services
[![Go Reference](https://pkg.go.dev/badge/github.com/mongodb-labs/go-client-mongodb-atlas-app-services.svg)](https://pkg.go.dev/github.com/mongodb-labs/go-client-mongodb-atlas-app-services)

A Go HTTP client for the [MongoDB Atlas App Services Admin API](https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/).

## Usage

```go
import "github.com/mongodb-labs/go-client-mongodb-atlas-app-services/appservices"
```

Construct a new App Services client, then use the various services on the client to
access different parts of the Atlas API. For example:

```go
client := appservices.NewClient(nil)
```

The services of a client divide the API into logical chunks and correspond to
the structure of the Atlas API documentation at
https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/.

**NOTE:** Using the [context](https://godoc.org/context) package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

## Versioning

Each version of the client is tagged, and the version is updated accordingly.

To see the list of past versions, run `git tag`.

To release a new version, first ensure that [Version](./appservices/appservices.go) is updated 
(i.e., before running `git push origin vx.y.z`, verify that `Version=x.y.z` should match the tag being pushed to GitHub)

## Roadmap

This library is being initially developed for [Atlas Terraform Provider](https://github.com/mongodb/terraform-provider-mongodbatlas)
so API methods will likely be implemented in the order that they are
needed by those projects.

## Contributing

See our [CONTRIBUTING.md](CONTRIBUTING.md) Guide.

## License

`go-client-mongodb-atlas-app-services` is released under the Apache 2.0 license. See [LICENSE](LICENSE)
