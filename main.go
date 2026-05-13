// main.go
package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"os"
	"database/sql"
	"gator/internal/config"
	"gator/internal/database"

)


func main(){
	// username := "natedawg"

	configFilePath := config.GetConfigFilePath()

	cfg, err := config.Read(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	s := state{db: dbQueries, cfg: &cfg}

	cmds := &commands{}
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerListUsers)
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", handlerAddFeed)
	cmds.Register("feeds", handlerListFeeds)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("usage: %s <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	 cmd := command{
		name: args[0],
		args: args[1:],
	}
	if err := cmds.Run(&s, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

}