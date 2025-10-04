/*
Copyright © 2025 NAME HERE dai.tsuruga0809@gmail.com
*/
package main

import (
	"github.com/Daaaai0809/swagen-v2/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	cmd.Execute()
}
