package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func rollDiceHandler(c echo.Context) error {
	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)
	diceRoll := rnd.Intn(6) + 1

	return c.JSON(http.StatusOK, map[string]int{"result": diceRoll})
}

func main() {
	e := echo.New()
	e.GET("/roll", rollDiceHandler)
	e.Logger.Fatal(e.Start(":8083"))
}
