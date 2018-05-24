package main

import (
	"flag"
	"github.com/sethvargo/go-password/password"

	log "github.com/sirupsen/logrus"

	. "git.codecoop.org/systemli/ticker/internal/api"
	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
	"git.codecoop.org/systemli/ticker/internal/bridge"
)

func main() {
	API().Run(Config.Listen)
}

func init() {
	var cp = flag.String("config", "config.yml", "path to config.yml")
	flag.Parse()

	Config = LoadConfig(*cp)
	DB = OpenDB(Config.Database)

	if Config.TwitterEnabled() {
		bridge.Twitter = bridge.NewTwitterBridge(Config.Twitter.ConsumerKey, Config.Twitter.ConsumerSecret)
	}

	firstRun()

	log.Print("starting ticker at ", Config.Listen)

	lvl, err := log.ParseLevel(Config.LogLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)
}

func firstRun() {
	count, err := DB.Count(&User{})
	if err != nil {
		log.Fatal("error using database")
	}

	if count == 0 {
		pw, err := password.Generate(24, 3, 3, false, false)
		if err != nil {
			log.Fatal(err)
		}

		user, err := NewAdminUser(Config.Initiator, pw)
		if err != nil {
			log.Fatal("could not create first user")
		}

		err = DB.Save(user)
		if err != nil {
			log.Fatal("could not persist first user")
		}

		log.WithField("email", user.Email).WithField("password", pw).Info("admin user created (change password now!)")
	}
}
