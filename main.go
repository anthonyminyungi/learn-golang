package main

import (
	"os"
	"strings"

	"github.com/anthonyminyungi/learngo/scrapper"
	"github.com/labstack/echo/v4"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

var fileName string = "jobs.csv"

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(fileName, term+"Jobs.csv")
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
