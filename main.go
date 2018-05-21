package main

import (
	"log"

	"github.com/sethvargo/go-password/password"

	. "git.codecoop.org/systemli/ticker/internal/api"
	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

func main() {
	API().Run(":8080")
}

func init() {
	//TODO: Config/Flags for database and ports
	DB = OpenDB("ticker.db")

	count, err := DB.Count(&User{})
	if err != nil {
		log.Fatal("Error using database")
	}

	if count == 0 {
		pw, err := password.Generate(24, 3, 3, false, false)
		if err != nil {
			log.Fatal(err)
		}

		//TODO: Make Email configurable
		user, err := NewUser("admin@systemli.org", pw)
		if err != nil {
			log.Fatal("Could not create first user")
		}
		user.IsSuperAdmin = true

		err = DB.Save(user)
		if err != nil {
			log.Fatal("Could not persist first user")
		}

		log.Println("First run: Creating User")
		log.Println("=======================================")
		log.Printf("Password: %s\n", pw)
		log.Println("Please change the password immediately!")
		log.Println("=======================================")
	}
}
