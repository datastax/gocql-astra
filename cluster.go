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

	"github.com/apache/cassandra-gocql-driver/v2"
)

func NewClusterFromBundle(path, username, password string, timeout time.Duration) (*gocql.ClusterConfig, error) {
	return NewClusterFromBundleWithLogger(path, username, password, timeout, nil)
}

func NewClusterFromURL(url, databaseID, token string, timeout time.Duration) (*gocql.ClusterConfig, error) {
	return NewClusterFromURLWithLogger(url, databaseID, token, timeout, nil)
}

func NewCluster(dialer gocql.HostDialer, username, password string) *gocql.ClusterConfig {
	return NewClusterWithLogger(dialer, username, password, nil)
}

func NewClusterFromBundleWithLogger(path, username, password string, timeout time.Duration, logger gocql.StructuredLogger) (*gocql.ClusterConfig, error) {
	dialer, err := NewDialerFromBundleWithLogger(path, timeout, logger)
	if err != nil {
		return nil, err
	}
	return NewClusterWithLogger(dialer, username, password, logger), nil
}

func NewClusterFromURLWithLogger(url, databaseID, token string, timeout time.Duration, logger gocql.StructuredLogger) (*gocql.ClusterConfig, error) {
	dialer, err := NewDialerFromURLWithLogger(url, databaseID, token, timeout, logger)
	if err != nil {
		return nil, err
	}
	return NewClusterWithLogger(dialer, "token", token, logger), nil
}

func NewClusterWithLogger(dialer gocql.HostDialer, username, password string, logger gocql.StructuredLogger) *gocql.ClusterConfig {
	// add multiple fake contact points to make gocql call the dialer multiple times (since the dialer will cycle through the contact points
	cluster := gocql.NewCluster("0.0.0.1", "0.0.0.2", "0.0.0.3") // Placeholder, maybe figure how to make this better
	cluster.HostDialer = dialer

	cluster.PoolConfig = gocql.PoolConfig{HostSelectionPolicy: gocql.RoundRobinHostPolicy()}
	cluster.Authenticator = &gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.ReconnectInterval = 30 * time.Second
	if logger != nil {
		cluster.Logger = logger
	}
	return cluster
}
