package main

import "flag"

var (
    configPath = flag.String("config_path", "web.config.json", "")
)

func main()  {
    flag.Parse()

}
