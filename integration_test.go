package gocqlastra

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	astrasdk "github.com/datastax/astra-client-go/v2/astra"
	"github.com/datastax/cql-proxy/astra"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/require"
)

var flagBundle = flag.String("bundle", "", "path to the Astra secure connect bundle")
var flagToken = flag.String("token", "", "astra db token")
var flagUsername = flag.String("username", "", "astra db username")
var flagPassword = flag.String("password", "", "astra db password")
var flagDbName = flag.String("db_name", "", "astra db name")
var flagApiUrl = flag.String("api_url", AstraAPIURL, "astra api url")

var dbId string

func TestMain(m *testing.M) {
	flag.Parse()

	if *flagBundle == "" {
		_, _ = fmt.Fprintln(os.Stderr, "-bundle is required")
		os.Exit(1)
	}

	if *flagToken == "" {
		_, _ = fmt.Fprintln(os.Stderr, "-token is required")
		os.Exit(1)
	}

	if *flagUsername == "" {
		_, _ = fmt.Fprintln(os.Stderr, "-username is required")
		os.Exit(1)
	}

	if *flagPassword == "" {
		_, _ = fmt.Fprintln(os.Stderr, "-password is required")
		os.Exit(1)
	}

	id, err := getDbId()
	assertInit(err == nil, "%s", err)

	dbId = id

	code := m.Run() // Run all tests

	os.Exit(code)
}

func assertInit(condition bool, msg string, args ...any) {
	if !condition {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", msg), args...)
		os.Exit(1)
	}
}

func getDbId() (string, error) {
	client, err := astrasdk.NewClientWithResponses(*flagApiUrl, func(c *astrasdk.Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *flagToken))
			return nil
		})
		return nil
	})

	assertInit(err == nil, "Failed to create client: %s", err)
	ctx, fn := context.WithTimeout(context.Background(), 30*time.Second)
	defer fn()
	resp, err := client.ListDatabasesWithResponse(ctx, &astrasdk.ListDatabasesParams{})
	assertInit(err == nil, "Failed to get databases: %s", err)
	assertInit(resp.JSON200 != nil, "bad response from list databases call: %v", string(resp.Body))
	assertInit(len(*resp.JSON200) != 0, "0 databases returned")
	for _, db := range *resp.JSON200 {
		if db.Info.Name != nil && *db.Info.Name == *flagDbName {
			return db.Id, nil
		}
	}

	return "", fmt.Errorf("could not find database '%s'", *flagDbName)
}

func coreTest(t *testing.T, c *gocql.ClusterConfig) {
	session, err := c.CreateSession()
	require.Nil(t, err)
	result, err := session.Query("SELECT * FROM system.local").Iter().SliceMap()
	require.Nil(t, err)
	fmt.Println(result)
	result, err = session.Query("SELECT * FROM system.peers").Iter().SliceMap()
	require.Nil(t, err)
	fmt.Println(result)
}

func TestNewCluster(t *testing.T) {
	bundle, err := astra.LoadBundleZipFromPath(*flagBundle)
	d, err := NewDialer(bundle, 30*time.Second)
	require.Nil(t, err)
	c := NewCluster(d, *flagUsername, *flagPassword)
	coreTest(t, c)
}

func TestNewClusterFromBundle(t *testing.T) {
	c, err := NewClusterFromBundle(*flagBundle, *flagUsername, *flagPassword, 30*time.Second)
	require.Nil(t, err)
	coreTest(t, c)
}

func TestNewClusterFromURL(t *testing.T) {
	c, err := NewClusterFromURL(*flagApiUrl, dbId, *flagToken, 30*time.Second)
	require.Nil(t, err)
	coreTest(t, c)
}
