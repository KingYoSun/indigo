package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/meilisearch/meilisearch-go"

	"github.com/bluesky-social/indigo/search"
	"github.com/bluesky-social/indigo/util/cliutil"

	"github.com/bluesky-social/indigo/util/version"
	logging "github.com/ipfs/go-log"
	es "github.com/opensearch-project/opensearch-go/v2"
	cli "github.com/urfave/cli/v2"
)

var log = logging.Logger("palomar")

func init() {
	// control log level using, eg, LOG_LEVEL=debug
	logging.SetAllLoggers(logging.LevelInfo)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {

	app := cli.App{
		Name:    "palomar",
		Usage:   "search indexing and query service (using ES or OS)",
		Version: version.Version,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "elastic-cert-file",
			Usage:   "certificate file path",
			EnvVars: []string{"ES_CERT_FILE", "ELASTIC_CERT_FILE"},
		},
		&cli.StringFlag{
			Name:    "elastic-username",
			Usage:   "elasticsearch username",
			Value:   "admin",
			EnvVars: []string{"ES_USERNAME", "ELASTIC_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "elastic-password",
			Usage:   "elasticsearch password",
			Value:   "admin",
			EnvVars: []string{"ES_PASSWORD", "ELASTIC_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "search-engine-hosts",
			Usage:   "searchengine hosts (schema/host/port)",
			Value:   "http://localhost:7700",
			EnvVars: []string{"ES_HOSTS", "ELASTIC_HOSTS", "MEILI_HOSTS", "SEARCH_ENGINE_HOSTS"},
		},
		&cli.StringFlag{
			Name:    "post-index",
			Usage:   "Index for 'post' documents",
			Value:   "posts",
			EnvVars: []string{"ES_POST_INDEX", "MEILI_POST_INDEX", "POST_INDEX"},
		},
		&cli.StringFlag{
			Name:    "profile-index",
			Usage:   "Index for 'profile' documents",
			Value:   "profiles",
			EnvVars: []string{"ES_PROFILE_INDEX", "MEILI_PROFILE_INDEX", "PROFILE_INDEX"},
		},
		&cli.StringFlag{
			Name:    "atp-bgs-host",
			Usage:   "hostname and port of BGS to subscribe to",
			Value:   "wss://bsky.social",
			EnvVars: []string{"ATP_BGS_HOST"},
		},
		&cli.StringFlag{
			Name:    "atp-plc-host",
			Usage:   "method, hostname, and port of PLC registry",
			Value:   "https://plc.directory",
			EnvVars: []string{"ATP_PLC_HOST"},
		},
		// TODO(bnewbold): this is a temporary hack to fetch our own blobs
		&cli.StringFlag{
			Name:    "atp-pds-host",
			Usage:   "method, hostname, and port of PDS instance",
			Value:   "https://bsky.social",
			EnvVars: []string{"ATP_PDS_HOST"},
		},
		&cli.StringFlag{
			Name:    "search-engine-type",
			Usage:   "elastic or meili",
			Value:   "meili",
			EnvVars: []string{"SEARCH_ENGINE_TYPE"},
		},
		&cli.StringFlag{
			Name:    "meilisearch-apikey",
			Usage:   "apikey for meiliseach",
			Value:   "meili-master-key",
			EnvVars: []string{"MEILISEARCH_APIKEY"},
		},
		&cli.StringFlag{
			Name:    "admin-key",
			EnvVars: []string{"ADMIN_KEY"},
		},
		&cli.IntFlag{
			Name:    "max-metadb-connections",
			EnvVars: []string{"MAX_METADB_CONNECTIONS"},
			Value:   40,
		},
	}

	app.Commands = []*cli.Command{
		elasticCheckCmd,
		searchCmd,
		runCmd,
	}

	return app.Run(args)
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "combined indexing+query server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "database-url",
			// XXX: data/palomar/search.db
			Value:   "sqlite://data/thecloud.db",
			EnvVars: []string{"DATABASE_URL"},
		},
		&cli.BoolFlag{
			Name:    "readonly",
			EnvVars: []string{"READONLY"},
		},
		&cli.StringFlag{
			Name:    "bind",
			Usage:   "IP or address, and port, to listen on for HTTP APIs",
			Value:   ":3999",
			EnvVars: []string{"PALOMAR_BIND"},
		},
	},
	Action: func(cctx *cli.Context) error {
		db, err := cliutil.SetupDatabase(cctx.String("database-url"), cctx.Int("max-metadb-connections"))
		if err != nil {
			return err
		}

		var meilicli *meilisearch.Client
		var escli *es.Client
		var clierr error
		if cctx.String("search-engine-type") == "meili" {
			meilicli = meilisearch.NewClient(meilisearch.ClientConfig{
				Host: cctx.String("search-engine-hosts"),
				APIKey: cctx.String("meilisearch-apikey"),
			})
		} else {
			escli, clierr = createEsClient(cctx)
			if clierr != nil {
				return fmt.Errorf("failed to get elasticsearch: %w", err)
			}
		}

		srv, err := search.NewServer(
			db,
			escli,
			meilicli,
			cctx.String("search-engine-type"),
			cctx.String("atp-plc-host"),
			cctx.String("atp-pds-host"),
			cctx.String("atp-bgs-host"),
		)
		if err != nil {
			return err
		}

		if tok := cctx.String("admin-key"); tok != "" {
			if err := srv.CreateAdminToken(tok); err != nil {
				return fmt.Errorf("failed to set up admin token: %w", err)
			}
		}

		go func() {
			srv.RunAPI(cctx.String("bind"))
		}()

		if cctx.Bool("readonly") {
			select {}
		} else {
			ctx := context.TODO()
			if err := srv.RunIndexer(ctx); err != nil {
				return fmt.Errorf("failed to run indexer: %w", err)
			}
		}

		return nil
	},
}

var elasticCheckCmd = &cli.Command{
	Name: "elastic-check",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "elastic-cert",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.String("search-engine-type") == "meili" {
			fmt.Println("Use Meilisearch")
			meilicli := meilisearch.NewClient(meilisearch.ClientConfig{
				Host: cctx.String("search-engine-hosts"),
				APIKey: cctx.String("meilisearch-apikey"),
			})
			res := meilicli.IsHealthy()
			if !res {
				return fmt.Errorf("Meilisearch is not Healthy")
			}

			return nil
		}

		escli, err := createEsClient(cctx)
		if err != nil {
			return err
		}

		// NOTE: this extra info check is redundant; createEsClient() already made this call and logged results
		inf, err := escli.Info()
		if err != nil {
			return fmt.Errorf("failed to get info: %w", err)
		}

		fmt.Println(inf)
		return nil

	},
}

var searchCmd = &cli.Command{
	Name:  "search",
	Usage: "run a simple query against search index",
	Action: func(cctx *cli.Context) error {
		if cctx.String("search-engine-type") == "meili" {
			fmt.Println("Use Meilisearch, skip searchCmd")
			return nil
		}

		escli, err := createEsClient(cctx)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"text": cctx.Args().First(),
				},
			},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}

		// Perform the search request.
		res, err := escli.Search(
			escli.Search.WithContext(context.Background()),
			escli.Search.WithIndex(cctx.String("posts-index")),
			escli.Search.WithBody(&buf),
			escli.Search.WithTrackTotalHits(true),
			escli.Search.WithPretty(),
		)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		fmt.Println(res)
		return nil

	},
}

func createEsClient(cctx *cli.Context) (*es.Client, error) {

	addrs := []string{}
	if hosts := cctx.String("search-engine-hosts"); hosts != "" {
		addrs = strings.Split(hosts, ",")
	}

	certfi := cctx.String("elastic-cert-file")
	var cert []byte
	if certfi != "" {
		b, err := os.ReadFile(certfi)
		if err != nil {
			return nil, err
		}

		cert = b
	}

	cfg := es.Config{
		Addresses: addrs,
		Username:  cctx.String("elastic-username"),
		Password:  cctx.String("elastic-password"),

		CACert: cert,
	}

	escli, err := es.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to set up client: %w", err)
	}

	info, err := escli.Info()
	if err != nil {
		return nil, fmt.Errorf("cannot get escli info: %w", err)
	}
	defer info.Body.Close()
	log.Info(info)

	return escli, nil
}
