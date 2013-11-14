package main

import (
	"atlantis/manager/client"
)

func main() {
	cli := client.New()
	cli.Run()
}
