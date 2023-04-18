package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/KingYoSun/indigo/api"
	"github.com/KingYoSun/indigo/bgs"
	"github.com/KingYoSun/indigo/blobs"
	"github.com/KingYoSun/indigo/carstore"
	cliutil "github.com/KingYoSun/indigo/cmd/gosky/util"
	"github.com/KingYoSun/indigo/events"
	"github.com/KingYoSun/indigo/indexer"
	"github.com/KingYoSun/indigo/notifs"
	"github.com/KingYoSun/indigo/plc"
	"github.com/KingYoSun/indigo/repomgr"
	"github.com/KingYoSun/indigo/version"

	_ "net/http/pprof"

	_ "github.com/joho/godotenv/autoload"

	logging "github.com/ipfs/go-log"
	"github.com/meilisearch/meilisearch-go"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/plugin/opentelemetry/tracing"
)

var log = logging.Logger("bigsky")

func init() {
	// control log level using, eg, GOLOG_LOG_LEVEL=debug
	//logging.SetAllLoggers(logging.LevelDebug)
}

func main() {
	run(os.Args)
}

func run(args []string) {

	app := cli.App{
		Name:    "bigsky",
		Usage:   "atproto BGS/firehose daemon",
		Version: version.Version,
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "jaeger",
		},
		&cli.StringFlag{
			Name:    "db-url",
			Usage:   "database connection string for BGS database",
			Value:   "sqlite://./data/bigsky/bgs.sqlite",
			EnvVars: []string{"DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "carstore-db-url",
			Usage:   "database connection string for carstore database",
			Value:   "sqlite://./data/bigsky/carstore.sqlite",
			EnvVars: []string{"CARSTORE_DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "meilisearch-url",
			Usage:   "host url for meilisearch",
			Value:   "http://localhost:7700",
			EnvVars: []string{"MEILISEARCH_URL"},
		},
		&cli.StringFlag{
			Name:    "meilisearch-apikey",
			Usage:   "apikey for meilisearch",
			Value:   "meili-master-key",
			EnvVars: []string{"MEILISEARCH_APIKEY"},
		},
		&cli.BoolFlag{
			Name: "db-tracing",
		},
		&cli.StringFlag{
			Name:    "data-dir",
			Usage:   "path of directory for CAR files and other data",
			Value:   "data/bigsky",
			EnvVars: []string{"DATA_DIR"},
		},
		&cli.StringFlag{
			Name:    "plc-host",
			Usage:   "method, hostname, and port of PLC registry",
			Value:   "https://plc.directory",
			EnvVars: []string{"ATP_PLC_HOST"},
		},
		&cli.BoolFlag{
			Name:  "crawl-insecure-ws",
			Usage: "when connecting to PDS instances, use ws:// instead of wss://",
		},
		&cli.BoolFlag{
			Name:  "aggregation",
			Value: true,
		},
		&cli.StringFlag{
			Name:  "api-listen",
			Value: ":2470",
		},
		&cli.StringFlag{
			Name:  "debug-listen",
			Value: "localhost:2471",
		},
		&cli.StringFlag{
			Name: "disk-blob-store",
		},
		&cli.StringFlag{
			Name: "admin-pass",
			Value: "DevAdminPass",
			EnvVars: []string{"ADMIN_PASS"},
		},
	}

	app.Action = func(cctx *cli.Context) error {

		if cctx.Bool("jaeger") {
			url := "http://localhost:14268/api/traces"
			exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
			if err != nil {
				return err
			}
			tp := tracesdk.NewTracerProvider(
				// Always be sure to batch in production.
				tracesdk.WithBatcher(exp),
				// Record information about this application in a Resource.
				tracesdk.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("bgs"),
					attribute.String("environment", "test"),
					attribute.Int64("ID", 1),
				)),
			)

			otel.SetTracerProvider(tp)
		}

		// ensure data directory exists; won't error if it does
		datadir := cctx.String("data-dir")
		csdir := filepath.Join(datadir, "carstore")
		if err := os.MkdirAll(datadir, os.ModePerm); err != nil {
			return err
		}

		dburl := cctx.String("db-url")
		db, err := cliutil.SetupDatabase(dburl)
		if err != nil {
			return err
		}

		csdburl := cctx.String("carstore-db-url")
		csdb, err := cliutil.SetupDatabase(csdburl)
		if err != nil {
			return err
		}

		meilicli := meilisearch.NewClient(meilisearch.ClientConfig{
			Host: cctx.String("meilisearch-url"),
			APIKey: cctx.String("meilisearch-apikey"),
		})

		if cctx.Bool("db-tracing") {
			if err := db.Use(tracing.NewPlugin()); err != nil {
				return err
			}
			if err := csdb.Use(tracing.NewPlugin()); err != nil {
				return err
			}
		}

		os.MkdirAll(filepath.Dir(csdir), os.ModePerm)
		cstore, err := carstore.NewCarStore(csdb, csdir)
		if err != nil {
			return err
		}

		didr := &api.PLCServer{Host: cctx.String("plc-host")}
		cachedidr := plc.NewCachingDidResolver(didr, time.Minute*5, 1000)

		kmgr := indexer.NewKeyManager(cachedidr, nil)

		repoman := repomgr.NewRepoManager(db, cstore, kmgr)

		dbp, err := events.NewDbPersistence(db, cstore)
		if err != nil {
			return fmt.Errorf("setting up db event persistence: %w", err)
		}

		evtman := events.NewEventManager(dbp)

		go evtman.Run()

		notifman := &notifs.NullNotifs{}

		ix, err := indexer.NewIndexer(db, meilicli, notifman, evtman, cachedidr, repoman, true, cctx.Bool("aggregation"))
		if err != nil {
			return err
		}

		repoman.SetEventHandler(func(ctx context.Context, evt *repomgr.RepoEvent) {
			if err := ix.HandleRepoEvent(ctx, evt); err != nil {
				log.Errorw("failed to handle repo event", "err", err)
			}
		})

		var blobstore blobs.BlobStore
		if bsdir := cctx.String("disk-blob-store"); bsdir != "" {
			blobstore = &blobs.DiskBlobStore{bsdir}
		}

		bgs, err := bgs.NewBGS(db, ix, meilicli, repoman, evtman, cachedidr, blobstore, !cctx.Bool("crawl-insecure-ws"), cctx.String("admin-pass"))
		if err != nil {
			return err
		}

		// set up pprof endpoint
		go func() {
			if err := bgs.StartDebug(cctx.String("debug-listen")); err != nil {
				panic(err)
			}
		}()

		return bgs.Start(cctx.String("api-listen"))
	}

	app.RunAndExitOnError()
}
