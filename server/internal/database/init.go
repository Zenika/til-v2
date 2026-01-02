package database

import (
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zenika/tilv2back/internal/configuration"
	"os"
	"regexp"
)

//go:embed migrations/*
var migrationFolder embed.FS

var Database *sql.DB

func init() {
	openConnection()

	makeMigrations()
}

// openConnection gets a connection to the database and stores it in Database variable.
func openConnection() {
	var err error

	configuration.Logger.Info("Trying to connect to database...")

	Database, err = sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", configuration.Configuration.DatabaseFileName))
	if err != nil {
		configuration.Logger.Error(err.Error())
		os.Exit(1)
	}

	Database.SetMaxOpenConns(1)

	// Enable foreign keys as sqlite disables it by default on new connection!
	_, _ = Database.Exec("PRAGMA foreign_keys = ON;")

	configuration.Logger.Info("Connection established to database")
}

// makeMigrations is in charge of migrationFolder in our database. It takes each files in "migrationFolder" folder and applies them by ascending file name.
// If a migration was already made, it will not run it.
func makeMigrations() {

	configuration.Logger.Info("Applying migrations if needed...")

	if _, err := Database.Exec("CREATE TABLE IF NOT EXISTS `migrations`(`file_name` TEXT NOT NULL PRIMARY KEY, `migration_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL);"); err != nil {
		configuration.Logger.Error("Error occurred while creating migrations history table.", err)
		os.Exit(1)
	}

	files, _ := migrationFolder.ReadDir("migrations")
	for _, k := range files {
		// Check if file match name format for migrations
		if result, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*.sql$", k.Name()); !result {
			configuration.Logger.Error(fmt.Sprintf("File %s does not match YYYY-MM-DD-(.*).sql format, ignoring it!", k.Name()))
			continue
		}

		if exists, err := checkIfMigrationExists(k.Name()); err != nil {
			configuration.Logger.Error("Error occurred while verifying migration! Program will now exit to avoid data loss and/or corruption.", err)
			os.Exit(1)
		} else if !exists {
			configuration.Logger.Info(fmt.Sprintf("Applying migration %s...", k.Name()))
			migrationFile, err := migrationFolder.ReadFile(fmt.Sprintf("migrations/%s", k.Name()))

			if err != nil {
				configuration.Logger.Error("Error occurred while applying migration file when reading content. Program will now exit to avoid data loss and/or corruption.", err)
				os.Exit(1)
			}

			if _, err = Database.Exec("BEGIN TRANSACTION;"); err != nil {
				configuration.Logger.Error("Error occurred while applying migration file while creating transaction. Program will now exit to avoid data loss and/or corruption.", err)
				os.Exit(1)
			}

			if _, err = Database.Exec(string(migrationFile)); err != nil {
				_, _ = Database.Exec("ROLLBACK TRANSACTION;") // Rollback transaction but don't check for errors. As the program is going to exit, transaction will be lost anyway. That's just a stub to be clean.
				configuration.Logger.Error("Error occurred while applying migration file. Program will now exit to avoid data loss and/or corruption.", err)
				os.Exit(1)
			}

			if _, err = Database.Exec("COMMIT TRANSACTION;"); err != nil {
				configuration.Logger.Error("Error occurred while commiting migration file. Program will now exit to avoid data loss and/or corruption.", err)
				os.Exit(1)
			}

			// Adding migration to the table to avoid running it again
			if _, err = Database.Exec("INSERT INTO `migrations`(`file_name`) VALUES (?);", k.Name()); err != nil {
				configuration.Logger.Error(fmt.Sprintf("Unable to add migration file %s to database! Please add it manually, otherwise this migration will be applied again on the next start!", k.Name()), err)
			}

			configuration.Logger.Info(fmt.Sprintf("Migration %s applied successfully!", k.Name()))

		} else {
			configuration.Logger.Debug(fmt.Sprintf("Migration %s already applied, skipping...", k.Name()))
		}
	}

	configuration.Logger.Info("Migrations applied (if needed), database is now up to date.")
}

func checkIfMigrationExists(name string) (bool, error) {
	rows, err := Database.Query("SELECT * FROM `migrations` WHERE `file_name`=? LIMIT 1;", name)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
