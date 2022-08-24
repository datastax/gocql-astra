package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/kong"
	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
)

type runConfig struct {
	AstraBundle     string        `yaml:"astra-bundle" help:"Path to secure connect bundle for an Astra database. Requires '--username' and '--password'. Ignored if using the token or contact points option." short:"b" env:"ASTRA_BUNDLE"`
	AstraToken      string        `yaml:"astra-token" help:"Token used to authenticate to an Astra database. Requires '--astra-database-id'. Ignored if using the bundle path or contact points option." short:"t" env:"ASTRA_TOKEN"`
	AstraDatabaseID string        `yaml:"astra-database-id" help:"Database ID of the Astra database. Requires '--astra-token'" short:"i" env:"ASTRA_DATABASE_ID"`
	AstraApiURL     string        `yaml:"astra-api-url" help:"URL for the Astra API" default:"https://api.astra.datastax.com" env:"ASTRA_API_URL"`
	AstraTimeout    time.Duration `yaml:"astra-timeout" help:"Timeout for contacting Astra when retrieving the bundle and metadata" default:"10s" env:"ASTRA_TIMEOUT"`
	Username        string        `yaml:"username" help:"Username to use for authentication" short:"u" env:"USERNAME"`
	Password        string        `yaml:"password" help:"Password to use for authentication" short:"p" env:"PASSWORD"`
}

func main() {
	var cfg runConfig

	parser, err := kong.New(&cfg)
	if err != nil {
		panic(err)
	}

	var cliCtx *kong.Context
	if cliCtx, err = parser.Parse(os.Args[1:]); err != nil {
		parser.Fatalf("error parsing flags: %v", err)
	}

	var cluster *gocql.ClusterConfig
	if len(cfg.AstraBundle) > 0 {
		cluster, err = gocqlastra.NewClusterFromBundle(cfg.AstraBundle, cfg.Username, cfg.Password, cfg.AstraTimeout)
		if err != nil {
			cliCtx.Fatalf("unable to open bundle %s from file: %v", cfg.AstraBundle, err)
		}
	} else if len(cfg.AstraToken) > 0 {
		if len(cfg.AstraDatabaseID) == 0 {
			cliCtx.Fatalf("database ID is required when using a token")
		}
		cluster, err = gocqlastra.NewClusterFromURL(cfg.AstraApiURL, cfg.AstraDatabaseID, cfg.AstraToken, cfg.AstraTimeout)
		if err != nil {
			cliCtx.Fatalf("unable to load bundle %s from astra: %v", cfg.AstraBundle, err)
		}
	} else {
		cliCtx.Fatalf("must provide either bundle path or token")
	}

	start := time.Now()
	session, err := gocql.NewSession(*cluster)
	elapsed := time.Now().Sub(start)
	if err != nil {
		log.Fatalf("unable to connect session: %v", err)
	}

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
