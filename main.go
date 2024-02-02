// This is custom goose binary with sqlite3 support only.

package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()
	log.Default().Printf("%v", args)
	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	log.Default().Printf("%v", command)

	dbstring, exists := os.LookupEnv("DSN")

	if !exists {
		log.Fatal("DB Connection string not set")
	}

	db, err := goose.OpenDBWithDriver("postgres", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}
}
