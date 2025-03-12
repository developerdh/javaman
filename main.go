package main

import (
	"fmt"
	"javaman/cmd"
	"os"
)

var version = "dev"

func main() {
	// 处理版本标志
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("javaman version %s\n", version)
		os.Exit(0)
	}

	cmd.Execute()
}
