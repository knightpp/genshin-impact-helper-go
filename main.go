package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"helper/account"
	"io"
	"log"
	"os"
	"time"
)

func loadState(path string) time.Time {
	file, err := os.Open(path)
	if err != nil {
		log.Print("Error reading state file: ", err)
		return time.Unix(0, 0)
	}
	defer file.Close()
	var t time.Time
	err = json.NewDecoder(file).Decode(&t)
	if err != nil {
		log.Print("Error decoding state json: ", err)
		return time.Unix(0, 0)
	}
	return t
}

func saveState(path string, t time.Time) {
	file, err := os.Create(path)
	if err != nil {
		log.Print("Error creating state file: ", err)
		return
	}
	defer file.Close()
	bts, err := json.Marshal(&t)
	if err != nil {
		log.Print("Error encoding state json: ", err)
		return
	}
	io.Copy(file, bytes.NewBuffer(bts))
}

func main() {
	path := flag.String("file", "cookie.txt", "path to cookie file")
	statePath := flag.String("state", "state.json", "where to save/load state data")
	flag.Parse()
	log.Println("Started")
	beijing, err := time.LoadLocation("Asia/Shanghai") // "China/Beijing"
	if err != nil {
		log.Fatal(err)
	}
	lastSignIn := loadState(*statePath)
	nextDayAfterSignIn := lastSignIn.Add(time.Duration(24-lastSignIn.Hour()) * time.Hour)

	today := time.Now().In(beijing)
	if nextDayAfterSignIn.After(today) {
		log.Print("The site should reset at ~", nextDayAfterSignIn,
			" but now ", today)
		return
	}

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
	saveState(*statePath, time.Now().In(beijing))
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
