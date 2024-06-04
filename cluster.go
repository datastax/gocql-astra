// Copyright (c) DataStax, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocqlastra

import (
	"time"

	"github.com/gocql/gocql"
)

const apacheAuthenticator = "org.apache.cassandra.auth.PasswordAuthenticator"
const dseAuthenticator = "com.datastax.bdp.cassandra.auth.DseAuthenticator"
const astraAuthenticator = "org.apache.cassandra.auth.AstraAuthenticator"

func NewClusterFromBundle(path, username, password string, timeout time.Duration, coordinatorIdx int) (*gocql.ClusterConfig, error) {
	dialer, err := NewDialerFromBundle(path, timeout, coordinatorIdx)
	if err != nil {
		return nil, err
	}
	return NewCluster(dialer, username, password), nil
}

func NewClusterFromURL(url, databaseID, token string, timeout time.Duration) (*gocql.ClusterConfig, error) {
	dialer, err := NewDialerFromURL(url, databaseID, token, timeout)
	if err != nil {
		return nil, err
	}
	return NewCluster(dialer, "token", token), nil
}

func NewCluster(dialer gocql.HostDialer, username, password string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster("127.0.0.1") // Placeholder, maybe figure how to make this better
	cluster.HostDialer = dialer
	cluster.PoolConfig = gocql.PoolConfig{HostSelectionPolicy: gocql.RoundRobinHostPolicy()}
	cluster.Authenticator = &gocql.PasswordAuthenticator{
		Username:              username,
		Password:              password,
		AllowedAuthenticators: []string{apacheAuthenticator, dseAuthenticator, astraAuthenticator},
	}
	return cluster
}
