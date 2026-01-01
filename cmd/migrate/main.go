package main

import (
	"fmt"
	"log"
	"os"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/leviszaboo/go-api/config"
	"github.com/leviszaboo/go-api/db"
)

func main() {
	db, err := db.NewMySQLStorage(mysqlCfg.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		MultiStatements:      true,
	})
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "version" {
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("current version: %d, dirty: %v", version, dirty)
		return
	}

	if cmd == "force" {
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate force <version>")
		}
		var version int
		if _, err := fmt.Sscanf(os.Args[1], "%d", &version); err != nil {
			log.Fatal("invalid version number")
		}
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
		log.Printf("forced version to %d", version)
		return
	}

	if cmd == "up" {
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("no new migrations to apply")
			} else {
				log.Fatal(err)
			}
		} else {
			log.Println("migrations applied successfully")
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("no migrations to rollback")
			} else {
				log.Fatal(err)
			}
		} else {
			log.Println("migrations rolled back successfully")
		}
	}
}
