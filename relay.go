package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	scribble "github.com/nanobox-io/golang-scribble"
)

var db *scribble.Driver

func main() {
	initDatabase()

	e := echo.New()
	e.HideBanner = true

	e.GET("/", indexHandle)
	// http.HandleFunc("/register", registerCredential)
	//Handle start registration
	//Handle finish registration
	//Handle start authentication
	// Handle finish authentication

	fmt.Println("Starting server on port :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func initDatabase() {
	var err error
	db, err = scribble.New("data", &scribble.Options{})
	if err != nil {
		panic(err)
	}
}

func registerCredential(w http.ResponseWriter, r *http.Request) {
	//
}

type htmlTemplateData struct {
	Title string
}

func indexHandle(c echo.Context) error {
	return c.File("index.html")
}
