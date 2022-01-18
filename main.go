package main

import (
	"context"
	"helper/account"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var log = zap.S()

func handler(rw http.ResponseWriter, r *http.Request) {
	log := log.With("user-agent", r.UserAgent(), "remote-addr", r.RemoteAddr)
	errorInternalServer := func(err error, msg string) {
		log.With(err).Error(msg)
		status := http.StatusInternalServerError
		http.Error(rw, http.StatusText(status), status)
	}
	errorBadRequest := func(err error, msg string) {
		log.With(err).Error(msg)
		status := http.StatusBadRequest
		http.Error(rw, http.StatusText(status), status)
	}

	log.Debug("new connection")
	if r.Method != http.MethodGet {
		errorBadRequest(nil, "unexpected method: "+r.Method)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorBadRequest(err, "parse form failed")
		return
	}

	cookie := r.FormValue("cookie")
	if cookie == "" {
		errorBadRequest(nil, "empty cookie")
		return
	}

	acc := account.New(cookie)
	info, err := acc.GetInfo()
	if err != nil {
		errorInternalServer(err, "GetInfo failed")
		return
	}

	if !info.Data.IsSign {
		log.Info("going to sign-in")
		err = acc.SignIn()
		if err != nil {
			errorInternalServer(err, "sign-in failed")
			return
		}
		rw.WriteHeader(http.StatusOK)
	} else {
		log.Info("you've already signed-in")
		rw.WriteHeader(http.StatusNotModified)
	}
}

func main() {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zap.DebugLevel)
	config.DisableStacktrace = true
	logger, err := config.Build()
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("couldn't build logger")
	}
	log = logger.Sugar()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("enviromental variable PORT is empty")
		return
	}

	log.Infof("running server on %s port", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/check-in", handler)
	srv := http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		err := srv.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
		default:
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatal("couldn't shutdown properly", err)
	}
}
