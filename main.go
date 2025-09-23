package main

import (
	"context"
	"log"
	"os"

	"github.com/xescugc/texture-uv/cmd"
)

func main() {
	err := cmd.Cmd.Run(context.TODO(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
