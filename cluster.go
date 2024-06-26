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
	"net"
	"time"

	"github.com/gocql/gocql"
)


const apacheAuthenticator = "org.apache.cassandra.auth.PasswordAuthenticator"
const dseAuthenticator = "com.datastax.bdp.cassandra.auth.DseAuthenticator"
const astraAuthenticator = "org.apache.cassandra.auth.AstraAuthenticator"

func NewClusterFromBundle(path, username, password string, timeout time.Duration) (*gocql.ClusterConfig, error) {
	dialer, err := NewDialerFromBundle(path, timeout)
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
	// add multiple fake contact points to make gocql call the dialer multiple times (since the dialer will cycle through the contact points
	cluster := gocql.NewCluster("0.0.0.1", "0.0.0.2", "0.0.0.3") // Placeholder, maybe figure how to make this better
	cluster.HostDialer = dialer

	// this will make gocql ignore the contact point address for the control host initially and use the system.local address right away
	// while also preventing a panic in `ConnectAddress()` if the control connection fails to initialize
	cluster.AddressTranslator = gocql.AddressTranslatorFunc(func(addr net.IP, port int) (net.IP, int) {
		return net.IPv4zero, port
	})

	cluster.PoolConfig = gocql.PoolConfig{HostSelectionPolicy: gocql.RoundRobinHostPolicy()}
	cluster.Authenticator = &gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
		AllowedAuthenticators: []string{apacheAuthenticator, dseAuthenticator, astraAuthenticator},
	}
	cluster.ReconnectInterval = 30 * time.Second
	return cluster
}
