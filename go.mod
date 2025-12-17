module github.com/datastax/gocql-astra

go 1.19

// TODO: remove this when https://github.com/apache/cassandra-gocql-driver/pull/1920 is merged
replace github.com/apache/cassandra-gocql-driver/v2 v2.0.0 => github.com/worryg0d/gocql/v2 v2.0.0-20251217085152-d19ec6081932

require (
	github.com/apache/cassandra-gocql-driver/v2 v2.0.0
	github.com/datastax/cql-proxy v0.1.6
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/datastax/astra-client-go/v2 v2.2.54 // indirect
	github.com/datastax/go-cassandra-native-protocol v0.0.0-20220706104457-5e8aad05cf90 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.12.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
