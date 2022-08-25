# gocql for Astra (prototype)

This provides a custom `gocql.HostDialer` that can be used to allow gocql to connect to DataStax Astra. The goal is to
provide native support for gocql on Astra.

This library was made possible by the `gocql.HostDialer` interface added here: https://github.com/gocql/gocql/pull/1629

Note: Only works with a version of `gocql` with this [fix](https://github.com/gocql/gocql/commit/dc449c49ae76d903ee369128ccb296656643ab51).

Use this command to pull the correct version until the next release of `gocql`:

```
go get github.com/gocql/gocql@ce100a15a6899a7f42fbdc588874a36afcadc921
```

## Issues

* Astra uses Stargate which doesn't current support the system table `system.peers_v2`. Also, the underlying storage 
  system for Astra is returns `4.0.0.6816` for the `release_version` column, but it doesn't actually support Apache
  Cassandra 4.0 (which includes `system.peers_v2`). To work correctly it currently requires at least a version of 
  `gocql` with the following [fix](https://github.com/gocql/gocql/commit/dc449c49ae76d903ee369128ccb296656643ab51):
  * Here's the `gocql` PR to fix the issue: https://github.com/gocql/gocql/pull/1646
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
