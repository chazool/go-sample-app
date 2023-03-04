package main

import (
	_ "github.com/chazool/go-sample-app/common/docs"
	"github.com/chazool/go-sample-app/common/pkg/config"
)

func init() {
	config.InitConfig()
}

// @title Sample Service
// @version 1.0
// description Restful service for sample

// @contact.name Chazool
// @contact.url  https://lk.linkedin.com/in/chazool
// @contact.email chazoolk@gmail.com

//BasePath /api/v1
//@schemes https

func main() {
	config.Start()
}
