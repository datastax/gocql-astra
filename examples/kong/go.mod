module github.com/datastax/gocql-astra/v2/example/kong

go 1.19

replace github.com/datastax/gocql-astra/v2 => ../..

// TODO: remove this replace and bump to v2.0.1 when gocql 2.0.1 is released
replace github.com/apache/cassandra-gocql-driver/v2 v2.0.0 => github.com/apache/cassandra-gocql-driver/v2 v2.0.1-0.20260320161859-b86c662e14e2

require (
	github.com/alecthomas/kong v0.6.1
	github.com/apache/cassandra-gocql-driver/v2 v2.0.0
	github.com/datastax/gocql-astra/v2 v2.0.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/datastax/astra-client-go/v2 v2.2.54 // indirect
	github.com/datastax/cql-proxy v0.1.6 // indirect
	github.com/datastax/go-cassandra-native-protocol v0.0.0-20220706104457-5e8aad05cf90 // indirect
	github.com/deepmap/oapi-codegen v1.12.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
