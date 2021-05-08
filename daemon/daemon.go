package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"helper/account"
	"io"
	"log"
	"os"
	"time"
)

type Daemon struct {
	cookiePath string
	statePath  string
}

func NewDaemon(cookiePath, statePath string) Daemon {
	return Daemon{cookiePath: cookiePath, statePath: statePath}
}

func (d Daemon) ReadCookie() (*account.Account, error) {
	cookieBytes, err := os.ReadFile(d.cookiePath)
	if err != nil {
		return nil, fmt.Errorf("error reading cookie file: %w", err)
	}
	cookie := string(cookieBytes)
	acc, err := account.New(cookie)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// Loops indefinitely, sleeps to the next day
func (d Daemon) Run(acc *account.Account) {
	const baseSleepDur time.Duration = 60 * time.Second
	errSleepDur := baseSleepDur
	beijing, err := time.LoadLocation("Asia/Shanghai") // "China/Beijing"
	if err != nil {
		log.Fatal("Couldn't load Beijing timezone: ", err)
	}
	lastSignIn := d.loadState()
	for {
		nextDayAfterSignIn := lastSignIn.Add(time.Duration(24-lastSignIn.Hour()) * time.Hour)

		today := time.Now().In(beijing)
		if nextDayAfterSignIn.After(today) {
			durToNextDay := nextDayAfterSignIn.Sub(today)
			log.Print("The site should reset in ~", durToNextDay,
				" and going to sleep for that duration")
			time.Sleep(durToNextDay + time.Minute)
			continue
		}
		info, err := acc.GetInfo()
		if err != nil {
			log.Print("GetInfo error: ", err)
			time.Sleep(errSleepDur)
			errSleepDur = errSleepDur + errSleepDur
			continue
		}
		log.Printf("GetInfo(): %+v", info)
		if !info.Data.IsSign {
			err = acc.SignIn()
			if err != nil {
				log.Print("Sign-in error: ", err)
			} else {
				log.Print("Successfully signed-in")
			}
		} else {
			log.Print("You have already signed-in today")
		}
		errSleepDur = baseSleepDur
		nowBeijing := time.Now().In(beijing)
		d.saveState(nowBeijing)
		lastSignIn = nowBeijing
	}
}

func (d Daemon) loadState() time.Time {
	file, err := os.Open(d.statePath)
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

func (d Daemon) saveState(t time.Time) {
	file, err := os.Create(d.statePath)
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
	_, err = io.Copy(file, bytes.NewBuffer(bts))
	if err != nil {
		log.Print("Error copying buffer: ", err)
	}
}
