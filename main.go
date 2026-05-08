package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"gator/internal/config"
	"os"

)


func main(){
	// username := "natedawg"
	configFilePath := config.GetConfigFilePath()

	cfg, err := config.Read(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	s := config.State{Cfg: &cfg}

	cmds := &config.Commands{}
	cmds.Register("login", config.HandlerLogin)
	
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("usage: %s <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	 cmd := config.Command{
		Name: args[0],
		Args: args[1:],
	}
	if err := cmds.Run(&s, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

}