package cmd

import (
	"github.com/asdine/storm"
	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/legacy"
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

			legacyStorage := legacy.NewLegacyStorage(oldDb)
			migration := legacy.NewMigration(legacyStorage, store)
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
