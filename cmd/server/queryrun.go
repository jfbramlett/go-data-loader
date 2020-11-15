package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// The serveCmd will execute the generate command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "queries the database",
	Run:   query,
}

func init() {
	// Here we create the command line flags for our app, and bind them to our package-local
	// config variable.
	flags := queryCmd.Flags()
	flags.Int("samplesize", 10, "the number of times to run the test")
	flags.Int("retrievesize", 100, "the number of assets to retrieve")
	flags.Bool("sample", false, "load the sample table or metadata table")
	flags.String("dsn", "root:password@tcp(localhost:3306)/test", "db connection string")

	// Add the "serve" sub-command to the root command.
	rootCmd.AddCommand(queryCmd)
}

func query(cmd *cobra.Command, args []string) {
	logCfg := logrus.New()
	logCfg.SetFormatter(&logrus.JSONFormatter{})
	logger := logCfg.WithField("service", "dbtest")
	ctx := context.WithValue(context.Background(), logKey, logger)

	samplesize, err := cmd.Flags().GetInt("samplesize")
	if err != nil {
		logger.Error(err, "failed to resolve samplesize")
		return
	}

	retrievesize, err := cmd.Flags().GetInt("retrievesize")
	if err != nil {
		logger.Error(err, "failed to resolve retrievesize")
		return
	}

	dsn, err := cmd.Flags().GetString("dsn")
	if err != nil {
		logger.Error(err, "failed to resolve records")
		return
	}

	loadSample, err := cmd.Flags().GetBool("sample")
	if err != nil {
		logger.Error(err, "failed to resolve sample flag")
		return
	}

	if loadSample {
		querySample(ctx, samplesize, retrievesize, dsn)
	} else {
		querymeta(ctx, samplesize, retrievesize, dsn)
	}
}

func querySample(ctx context.Context, samplesize int, retrievesize int, dsn string) {
	logger := ctx.Value(logKey).(*logrus.Entry)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.WithError(err).Errorf("Waiting for DB to be ready")
		return
	}

	executionTimes := make([]time.Duration, 0)

	for i := 0; i < samplesize; i++ {
		uuids, err := getUUIDs(ctx, db, "sample", retrievesize)
		if err != nil {
			return
		}

		result := []*Sample{}
		query, args, _ := sqlx.In(`select uuid, asset_uuid, chord, skey, bpm, name from sample where asset_uuid in (?)`, uuids)
		query = db.Rebind(query)

		start := time.Now()
		err = db.SelectContext(ctx, &result, query, args...)
		if err != nil {
			logger.WithError(err).Error("failed running query")
			return
		}
		executionTimes = append(executionTimes, time.Since(start))
	}

	avg := int64(0)
	for _, td := range executionTimes {
		avg = avg + td.Milliseconds()
		fmt.Printf("%vms\n", td.Milliseconds())
	}
	fmt.Printf("Avg response: %dms", avg/int64(len(executionTimes)))
}

func getUUIDs(ctx context.Context, db *sqlx.DB, table string, num int) ([]string, error) {
	logger := ctx.Value(logKey).(*logrus.Entry)
	logger.Info("retrieving uuids")

	var rows int
	err := db.GetContext(ctx, &rows, "select max(id) from "+table+" limit 1")
	if err != nil {
		logger.WithError(err).Error("failed getting record count")
		return nil, err
	}

	idSet := make(map[int]int)
	uuids := make(map[string]string)
	for {
		if len(uuids) == num {
			break
		}
		id := rand.Intn(rows)
		if _, found := idSet[id]; !found {
			var uuid string
			err = db.GetContext(ctx, &uuid, "select asset_uuid from "+table+" where id=? limit 1", id)
			if err != nil {
				logger.WithError(err).Error("failed getting asset_uuid")
				return nil, err
			}
			uuids[uuid] = uuid
		}
	}

	result := make([]string, 0)
	for k := range uuids {
		result = append(result, k)
	}
	logger.Info("retrieved uuids")
	return result, nil
}

func querymeta(ctx context.Context, samplesize int, retrievesize int, dsn string) {
	logger := ctx.Value(logKey).(*logrus.Entry)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.WithError(err).Errorf("Waiting for DB to be ready")
		return
	}

	executionTimes := make([]time.Duration, 0)

	for i := 0; i < samplesize; i++ {
		uuids, err := getUUIDs(ctx, db, "asset_data", retrievesize)
		if err != nil {
			return
		}

		logger.Infof("running query with uuids %v", uuids)
		result := []*Metadata{}
		query, args, _ := sqlx.In(`select ad.asset_uuid, ad.asset_metadata_id, md.name, md.datatype, ad.value 
										 from asset_data ad 
										 join metadata md on ad.asset_metadata_id = md.id 
										 where ad.asset_uuid in (?) order by ad.asset_uuid`, uuids)
		query = db.Rebind(query)

		start := time.Now()
		err = db.SelectContext(ctx, &result, query, args...)
		if err != nil {
			logger.WithError(err).Error("failed running query")
			return
		}
		took := time.Since(start)
		logger.Infof("query completed in %dms", took.Milliseconds())
		executionTimes = append(executionTimes, took)
	}

	avg := int64(0)
	for _, td := range executionTimes {
		avg = avg + td.Milliseconds()
		fmt.Printf("%vms\n", td.Milliseconds())
	}
	fmt.Printf("Avg response: %dms", avg/int64(len(executionTimes)))
}
