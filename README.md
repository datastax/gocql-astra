# gocql for Astra

This provides a custom `gocql.HostDialer` that can be used to allow gocql to connect to DataStax Astra. The goal is to
provide native support for gocql on Astra.

This library was made possible by the `gocql.HostDialer` interface added [here](https://github.com/gocql/gocql/pull/1629)
and released with version [1.2.1](https://github.com/gocql/gocql/releases/tag/v1.2.1).  As a result you _must_ use
version 1.2.1 or greater of gocql with gocql-astra.

## Issues

* Need to verify that topology/status events correctly update the driver when using Astra.
  * This seems to work correctly and was tested by removing Astra coordinators
* There is a bit of weirdness around contact points. I'm just using a place holder `"0.0.0.0"` (some valid IP address) 
  then the `HostDialer` provides a host ID from the metadata service when the host ID in the `HostInfo` is empty.

## How to use it:

Using an Astra bundle:

```go
cluster, err := gocqlastra.NewClusterFromBundle("/path/to/your/bundle.zip", 
	"<username>", "<password>", 10 * time.Second)

if err != nil {
    panic("unable to load the bundle")
}

session, err := gocql.NewSession(*cluster)

// ...
```

Using an Astra token:

```go
cluster, err = gocqlastra.NewClusterFromURL(gocqlastra.AstraAPIURL, 
	"<astra-database-id>", "<astra-token>", 10 * time.Second)

if err != nil {
panic("unable to load the bundle")
}

session, err := gocql.NewSession(*cluster)

// ...
```

Also, look at the [example](example) for more information.

### Running the example:

```
cd example
go build

# Using a bundle
./example --astra-bundle /path/to/bundle.zip --username <username> --password <password>

# Using a token
./example --astra-token <astra-token> --astra-database-id <astra-database-id> \
  [--astra-api-url <astra-api-url>]
```
