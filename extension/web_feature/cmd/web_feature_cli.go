package main

import "flag"

var (
	configPath = flag.String("config_path", "", "")
)

func main() {
	flag.Parse()
}
