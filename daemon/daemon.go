package daemon

import (
	"fmt"
	"helper/account"
	"helper/daemon/config"
	"log"
	"os"
	"time"
)

// Loops indefinitely, sleeps to the next day
func Run(accCfg *config.AccConfig, save chan bool) error {
	const baseSleepDur time.Duration = 60 * time.Second
	errSleepDur := baseSleepDur
	beijing, err := time.LoadLocation("Asia/Shanghai") // "China/Beijing"
	if err != nil {
		return fmt.Errorf("couldn't load Beijing timezone: %w", err)
	}
	acc := account.New(accCfg.Cookie)
	if err != nil {
		return fmt.Errorf("couldn't create account struct: %w", err)
	}
	log.Printf("Starting %s", accCfg.Name)
	log := log.New(os.Stderr, fmt.Sprintf("%s |", accCfg.Name),
		log.Ltime|log.Ldate|log.Lmsgprefix)
	lastSignIn := &accCfg.LastSignIn
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
		*lastSignIn = nowBeijing
		save <- true
	}
}
