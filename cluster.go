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
	return newCluster(dialer, username, password), nil
}

func NewClusterFromURL(url, databaseID, token string, timeout time.Duration) (*gocql.ClusterConfig, error) {
	dialer, err := NewDialerFromURL(url, databaseID, token, timeout)
	if err != nil {
		return nil, err
	}
	return newCluster(dialer, "token", token), nil
}

func newCluster(dialer gocql.HostDialer, username, password string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster("0.0.0.0") // Placeholder, maybe figure how to make this better
	cluster.HostFilter = &HostFilter{}
	cluster.HostDialer = dialer
	cluster.PoolConfig = gocql.PoolConfig{HostSelectionPolicy: gocql.RoundRobinHostPolicy()}
	cluster.Authenticator = &gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	return cluster
}
