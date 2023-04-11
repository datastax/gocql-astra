package main

import (
	"fmt"
	"log"
	"os"
	"time"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

func main() {

	var err error

	err = godotenv.Load()

	var cluster *gocql.ClusterConfig
	if len(os.Getenv("ASTRA_DB_SECURE_BUNDLE_PATH")) > 0 {
		cluster, err = gocqlastra.NewClusterFromBundle(os.Getenv("ASTRA_DB_SECURE_BUNDLE_PATH"), "token", os.Getenv("ASTRA_DB_APPLICATION_TOKEN"), 10*time.Second)
		if err != nil {
			err = fmt.Errorf("unable to open bundle %s from file: %v", os.Getenv("ASTRA_DB_SECURE_BUNDLE_PATH"), err)
			panic(err)
		}
	} else if len(os.Getenv("ASTRA_DB_APPLICATION_TOKEN")) > 0 {
		if len(os.Getenv("ASTRA_DB_ID")) == 0 {
			panic("database ID is required when using a token")
		}
		cluster, err = gocqlastra.NewClusterFromURL("https://api.astra.datastax.com", os.Getenv("ASTRA_DB_ID"), os.Getenv("ASTRA_DB_APPLICATION_TOKEN"), 10*time.Second)
		fmt.Println(cluster)
		if err != nil {
			fmt.Errorf("unable to load cluster %s from astra: %v", os.Getenv("ASTRA_DB_APPLICATION_TOKEN"), err)
		}
	} else {
		fmt.Errorf("must provide either bundle path or token")
	}

	start := time.Now()
	session, err := gocql.NewSession(*cluster)
	elapsed := time.Now().Sub(start)
	if err != nil {
		log.Fatalf("unable to connect session: %v", err)
	}

	fmt.Println("Making the query now")

	iter := session.Query("SELECT release_version FROM system.local").Iter()

	var version string
	for iter.Scan(&version) {
		fmt.Println(version)
	}

	if err = iter.Close(); err != nil {
		log.Printf("error running query: %v", err)
	}

	fmt.Printf("Connection process took %s\n", elapsed)
}

