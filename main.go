package main

import (
	. "git.codecoop.org/systemli/ticker/internal/api"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

func main() {
	//TODO: Config/Flags for database and ports
	DB = OpenDB("ticker.db")

	API().Run(":8080")
}
