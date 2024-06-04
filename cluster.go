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
	cluster := gocql.NewCluster("127.0.0.1", "127.0.0.2", "127.0.0.3") // Placeholder, maybe figure how to make this better
	cluster.HostDialer = dialer
	cluster.PoolConfig = gocql.PoolConfig{HostSelectionPolicy: gocql.RoundRobinHostPolicy()}
	cluster.Authenticator = &gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.ReconnectInterval = 15 * time.Second
	return cluster
}
