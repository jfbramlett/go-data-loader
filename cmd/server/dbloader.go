package main

import (
	"context"
	"math/rand"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type logkeyStruct struct{}

var (
	logKey = logkeyStruct{}
	chords = []string{"major", "minor"}
	keys   = []string{"a#", "c", "d", "d#", "e", "e#"}
	bpms   = []int{16, 32, 64, 120, 240}
)

// The serveCmd will execute the generate command
var dbloaderCmd = &cobra.Command{
	Use:   "dbload",
	Short: "loads the database",
	Run:   load,
}

func init() {
	// Here we create the command line flags for our app, and bind them to our package-local
	// config variable.
	flags := dbloaderCmd.Flags()
	flags.Int("records", 1000000, "the number of records to load")
	flags.Bool("sample", false, "load the sample table or metadata table")
	flags.String("dsn", "root:password@tcp(localhost:3306)/test", "db connection string")

	// Add the "serve" sub-command to the root command.
	rootCmd.AddCommand(dbloaderCmd)
}

func load(cmd *cobra.Command, args []string) {
	logCfg := logrus.New()
	logCfg.SetFormatter(&logrus.JSONFormatter{})
	logger := logCfg.WithField("service", "dbtest")
	ctx := context.WithValue(context.Background(), logKey, logger)

	records, err := cmd.Flags().GetInt("records")
	if err != nil {
		logger.Error(err, "failed to resolve records")
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
		loadSampleTable(ctx, records, dsn)
	} else {
		loadMetadataTable(ctx, records, dsn)
	}
}

type Sample struct {
	UUID      string `db:"uuid"`
	AssetUUID string `db:"asset_uuid"`
	Chord     string `db:"chord"`
	Key       string `db:"skey"`
	BPM       int    `db:"bpm"`
	Name      string `db:"name"`
}

func loadSampleTable(ctx context.Context, records int, dsn string) {
	logger := ctx.Value(logKey).(*logrus.Entry)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.WithError(err).Errorf("Waiting for DB to be ready")
		return
	}

	stmt := `insert into sample (uuid, asset_uuid, chord, skey, bpm, name) values (:uuid, :asset_uuid, :chord, :skey, :bpm, :name)`
	for i := 0; i < records; i++ {
		chordsRand := rand.Intn(len(chords))
		keysRand := rand.Intn(len(keys))
		bpmRand := rand.Intn(len(bpms))

		sampleUUID := uuid.New().String()
		s := &Sample{UUID: sampleUUID,
			AssetUUID: sampleUUID,
			Chord:     chords[chordsRand],
			Key:       keys[keysRand],
			BPM:       bpms[bpmRand],
			Name:      uuid.New().String(),
		}

		_, err := db.NamedExecContext(ctx, stmt, s)
		if err != nil {
			logger.WithError(err).Error("failed executing insert")
			return
		}

		if i%10 == 0 {
			logger.Infof("completed %d records", i)
		}
	}
}

type Metadata struct {
	AssetUUID         string `db:"asset_uuid"`
	AssetMetaDataId   int    `db:"asset_metadata_id"`
	AssetMetaName     string `db:"name"`
	AssetMetaDatatype string `db:"datatype"`
	Value             string `db:"value"`
}

func loadMetadataTable(ctx context.Context, records int, dsn string) {
	logger := ctx.Value(logKey).(*logrus.Entry)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logger.WithError(err).Errorf("Waiting for DB to be ready")
		return
	}

	stmt := `insert into asset_data (asset_uuid, asset_metadata_id, value) values (:asset_uuid, :asset_metadata_id, :value)`
	for i := 0; i < records; i++ {
		chordsRand := rand.Intn(len(chords))
		keysRand := rand.Intn(len(keys))
		bpmRand := rand.Intn(len(bpms))
		sampleUUID := uuid.New().String()

		_, err := db.NamedExecContext(ctx, stmt, &Metadata{
			AssetUUID:       sampleUUID,
			AssetMetaDataId: 1,
			Value:           chords[chordsRand],
		})
		if err != nil {
			logger.WithError(err).Error("failed executing insert")
			return
		}

		_, err = db.NamedExecContext(ctx, stmt, &Metadata{
			AssetUUID:       sampleUUID,
			AssetMetaDataId: 2,
			Value:           keys[keysRand],
		})
		if err != nil {
			logger.WithError(err).Error("failed executing insert")
			return
		}
		_, err = db.NamedExecContext(ctx, stmt, &Metadata{
			AssetUUID:       sampleUUID,
			AssetMetaDataId: 3,
			Value:           strconv.Itoa(bpms[bpmRand]),
		})
		if err != nil {
			logger.WithError(err).Error("failed executing insert")
			return
		}
		_, err = db.NamedExecContext(ctx, stmt, &Metadata{
			AssetUUID:       sampleUUID,
			AssetMetaDataId: 4,
			Value:           uuid.New().String(),
		})
		if err != nil {
			logger.WithError(err).Error("failed executing insert")
			return
		}

		if i%10 == 0 {
			logger.Infof("completed %d records", i)
		}
	}
}
