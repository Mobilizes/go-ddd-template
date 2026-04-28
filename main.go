package main

import (
	"fmt"
	"mob/ddd-template/cmd"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
		return
	}

	runInstruction := `
usage: go run main.go [options]
options:
  --serve
    On by default, will be overridden if other options are added
  --migrate
    Fresh migrate the database
  --seed
    Seed the existing database, will throw error if database hasn't been migrated
	`

	options := map[string]func(){
		"--serve":   cmd.Serve,
		"--migrate": cmd.Migrate,
		"--seed":    cmd.Seed,
	}

	cmds := []string{"--serve"}
	if len(os.Args) > 1 {
		cmds = []string{}
		for _, arg := range os.Args[1:] {
			if options[arg] == nil {
				panic(runInstruction)
			}
			if !slices.Contains(cmds, arg) {
				cmds = append(cmds, arg)
			}
		}
	}

	for _, cmd := range cmds {
		options[cmd]()
	}
}
