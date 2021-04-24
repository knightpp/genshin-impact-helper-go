package main

import (
	"fmt"
	"helper/account"
	"log"
	"os"
)

func main() {
	log.Println("Started")
	cookieBytes, err := os.ReadFile("cookie.txt")
	if err != nil {
		log.Fatalln(err)
	}
	cookie := string(cookieBytes)
	acc, err := account.New(cookie, "ru-ru")
	if err != nil {
		log.Fatalln(err)
	}
	// resp, err := acc.GetInfo()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Printf("%+v", resp)
	awards, err := acc.GetAwards()
	if err != nil {
		log.Fatalln(err)
	}
	totals := make(map[string]int)
	for _, a := range awards.Data.Awards {
		totals[a.Name] += a.Count
	}
	for name, count := range totals {
		fmt.Printf("%30v: %5v\n", name, count)
	}
}
