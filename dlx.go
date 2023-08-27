package main

import (
	"dlx/cmd"
	"dlx/internal/download/extractors/bilibili"
	"dlx/log"
	"fmt"
	"os"
)

func init() {
	log.InitSugar()
	bilibili.InitBilibili()
}

func main() {
	if log.Sugar != nil {
		defer log.Sugar.Sync()
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
