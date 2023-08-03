/*
Copyright © 2023 KAI CHU CHUNG
*/
package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/cage1016/alfred-xlsx2csv/cmd"
)

func main() {
	cmd.Execute()
}