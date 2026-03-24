# gocql for Astra

This provides a custom `gocql.HostDialer` that can be used to allow gocql to connect to DataStax Astra. The goal is to
provide native support for gocql on Astra.

This library relies on the following features of gocql:

* The ability to customize connection features via the [HostDialer interface](https://github.com/gocql/gocql/pull/1629)
* [Querying system.peers](https://github.com/gocql/gocql/pull/1646) if system.peers_v2 should be used but isn't available

You must use a version of gocql which supports both of these features. Use version >= 2.0.1 of the Apache Cassandra GoCQL Driver (`github.com/apache/cassandra-gocql-driver/v2`).

## Migration from v1 to v2

Version 2.0.0 of gocql-astra introduces breaking changes due to the migration of the underlying gocql driver to the Apache Software Foundation. The gocql project was donated to the ASF and, as part of version 2, changed its module path from `github.com/gocql/gocql` to `github.com/apache/cassandra-gocql-driver/v2`.

If you're upgrading from v1, follow these steps:

### 1. Update your go.mod

Change your module dependency from:
```go
require github.com/datastax/gocql-astra v1.x.x
```

To:
```go
require github.com/datastax/gocql-astra/v2 v2.0.0
```

### 2. Update your imports

Change your import statements from:
```go
import (
    gocqlastra "github.com/datastax/gocql-astra"
    "github.com/gocql/gocql"
)
```

To:
```go
import (
    gocqlastra "github.com/datastax/gocql-astra/v2"
    gocql "github.com/apache/cassandra-gocql-driver/v2"
)
```

### 3. Update your dependencies

Run the following commands to update your dependencies:
```bash
go get github.com/datastax/gocql-astra/v2@latest
go get github.com/apache/cassandra-gocql-driver/v2@latest
go mod tidy
```

### 4. Review API changes

The core API of gocql-astra remains the same, but you should review any code that directly uses the Apache Cassandra GoCQL Driver API, as it has undergone changes in v2. Refer to the [Apache Cassandra GoCQL Driver Upgrade Guide](https://github.com/apache/cassandra-gocql-driver/blob/trunk/UPGRADE_GUIDE.md) for details on driver-specific changes.


## Issues

* There is a bit of weirdness around contact points. The driver is using a few place holders `"0.0.0.1,0.0.0.2,0.0.0.3"` (some valid IP address) 
  then the `HostDialer` provides a host ID from the metadata service when the host ID in the `HostInfo` is empty. 
  Using multiple placeholder contact points instead of a single one enables the driver to retry if the initial connection fails.

## How to use it:

Using an Astra bundle:

```go
import (
	gocqlastra "github.com/datastax/gocql-astra/v2"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

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
import (
	gocqlastra "github.com/datastax/gocql-astra/v2"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

cluster, err = gocqlastra.NewClusterFromURL(gocqlastra.AstraAPIURL,
	"<astra-database-id>", "<astra-token>", 10 * time.Second)

if err != nil {
    panic("unable to load the bundle")
}

session, err := gocql.NewSession(*cluster)

// ...
```

Also, look at the [example](examples) for more information.

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
