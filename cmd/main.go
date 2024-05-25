package main

import (
	"time"

	"github.com/thisismz/data-processor/internal/app"
	_ "github.com/thisismz/data-processor/pkg/docs"
)

var (
	Version string
	Build   string
	Branch  string
)

// @title data-processor API
// @version 1.0
// @description This is a  data-processor server.

// @contact.name mahdi mozaffari
// @contact.email mahdimozaffari@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	if Version == "" {
		Version = "unknown"
	}
	if Branch == "" {
		Branch = "unknown"
	}
	if Build == "" {
		Build = time.Now().Format("2006-01-02 15:04:05")
	}
	app.Run()
}
