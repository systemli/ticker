package cmd

import (
	"github.com/asdine/storm"
	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/legacy"
	"github.com/systemli/ticker/internal/storage"
)

var (
	stormPath string

	dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Manage the database",
		Long:  "Commands for managing the database.",
		Args:  cobra.ExactArgs(1),
	}

	dbMigrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database from BoltDB to SQL",
		Run: func(cmd *cobra.Command, args []string) {
			oldDb, err := storm.Open(stormPath)
			if err != nil {
				log.WithError(err).Fatal("Unable to open the old database")
			}
			defer oldDb.Close()

			db, err := storage.OpenGormDB(cfg.Database.Type, cfg.Database.DSN, log)
			if err != nil {
				log.WithError(err).Fatal("Unable to open the new database")
			}

			if err := storage.MigrateDB(db); err != nil {
				log.WithError(err).Fatal("Unable to migrate the database")
			}

			legacyStorage := legacy.NewLegacyStorage(oldDb)
			newStorage := storage.NewSqlStorage(db, cfg.UploadPath)

			migration := legacy.NewMigration(legacyStorage, newStorage)
			if err := migration.Do(); err != nil {
				log.WithError(err).Fatal("Unable to migrate the database")
			}
		},
	}
)

func init() {
	dbCmd.AddCommand(dbMigrateCmd)

	dbMigrateCmd.Flags().StringVar(&stormPath, "storm.path", "", "path to the old db file")
}
