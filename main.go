// This is custom goose binary with sqlite3 support only.

package main

import (
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v\n", err)
	}
	args := flags.Args()
	log.Default().Printf("%v", args)
	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	log.Default().Printf("%v", command)

	dbstring := viper.GetString("dsn")

	if dbstring == "" {
		log.Fatalf("goose: missing dsn in config file\n")
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

	ctx := context.Background()

	if err := goose.RunContext(ctx, command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}

func init() {
	projectRoot := "."
	// Search config in home directory with name ".cobra.yaml".
	viper.AddConfigPath(projectRoot)
	viper.SetConfigType("yaml")
	viper.SetConfigName("goose.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("goose: please create a config file in the root of your project called goose.yaml\n")
	}
}
