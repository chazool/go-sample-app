package main

import "github.com/chazool/go-sample-app/common/pkg/config"

func init() {
	config.InitConfig()
}

func main() {
	config.Start()
}
