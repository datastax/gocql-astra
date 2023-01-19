# gocql for Astra

This provides a custom `gocql.HostDialer` that can be used to allow gocql to connect to DataStax Astra. The goal is to
provide native support for gocql on Astra.

This library relies on the following features of gocql:

* The ability to customize connection features via the [HostDialer interface](https://github.com/gocql/gocql/pull/1629)
* [Querying system.peers](https://github.com/gocql/gocql/pull/1646) if system.peers_v2 should be used but isn't available 

You must use a version of gocql which supports both of these features.  Both features have been merged into master as of
version [1.2.1](https://github.com/gocql/gocql/releases/tag/v1.2.1) so any release >= 1.2.1 should work.

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
