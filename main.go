package main

import (
	"flag"
	"helper/account"
	"log"
	"os"
)

func main() {
	path := flag.String("file", "cookie.txt", "path to cookie file")
	log.Println("Started")
	cookieBytes, err := os.ReadFile(*path)
	if err != nil {
		log.Fatalln(err)
	}
	cookie := string(cookieBytes)
	acc, err := account.New(cookie)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := acc.GetInfo()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("GetInfo(): %+v", info)
	if !info.Data.IsSign {
		err = acc.SignIn()
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Successfully signed-in\n")
	} else {
		log.Println("You have already signed-in today")
	}
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
