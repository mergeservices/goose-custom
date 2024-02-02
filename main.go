// This is custom goose binary with sqlite3 support only.

package main

import (
	"flag"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
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

	dbstring := viper.GetString("DSN")

	if dbstring == "" {
		log.Fatalf("goose: missing DSN")
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
	projectRoot := "."
	// Search config in home directory with name ".cobra.yaml".
	viper.AddConfigPath(projectRoot)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".cobra.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Println("goose: using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}
}
