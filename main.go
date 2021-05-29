package main

import (
	"flag"
	"helper/daemon"
	"helper/daemon/config"
	"log"
	"sync"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to a config.toml")
	flag.Parse()
	if configPath == "" {
		log.Fatal("config path is empty")
	}
	log.Println("Started")
	c, err := config.FromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	var save chan bool
	wg := sync.WaitGroup{}
	for i := range c.Account {
		wg.Add(1)
		go func(acc *config.AccConfig) {
			daemon.Run(acc, save)
			wg.Done()
		}(&c.Account[i])
	}
	go func() {
		for {
			<-save
			c.WriteToFile(configPath)
		}
	}()
	wg.Wait()
	log.Println("No more work to do, exitting...")
}

// func showTotals(acc *account.Account) {
// 	awards, err := acc.GetAwards()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	totals := make(map[string]int)
// 	for _, a := range awards.Data.Awards {
// 		totals[a.Name] += a.Count
// 	}
// 	for name, count := range totals {
// 		fmt.Printf("%30v: %5v\n", name, count)
// 	}
// }
