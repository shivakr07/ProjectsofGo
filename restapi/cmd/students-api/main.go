package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shivakr07/students-api/internal/config"
)

func main() {
	//load config
	// custom logger [if you have any]
	// db setup
	// router setup
	// server setup

	//load config
	cfg := config.MustLoad()

	//router setup
	//we will use net/http inbuilt package
	router := http.NewServeMux()
	//now we can make url's
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
		//in write method we can add bytes
	})

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// fmt.Printf("server started %s", cfg.HTTPServer.Addr)
	slog.Info("server started", slog.String("address", cfg.Addr))

	// but generally this is not how we keep our server
	// because if some interruption happens from user like C^Signal = interrupt or other then it will immediate terminate the server, but there might be some request which is in processing so first we need to complete that and then shutdown [called graceful shutdown]

	// err := server.ListenAndServe()
	// //this is blocking so keeping print before this
	// if err != nil {
	// 	log.Fatal("failed to start server")
	// }
	// so we keep our server in other go routine

	// --- SERVER WITH GRACEFUL SHUTDOWN ---
	// to synchronize the go routine we use wg/channel

	done := make(chan os.Signal, 1) //buffered

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// if any mentioned signal comes from os or user then notify in the channel

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done
	// it will be unblocked when channel receives the signal

	//now then we will shutdown...
	// structured logging
	slog.Info("shutting down the server")
	// server.Shutdown()
	// this gracefully shutdown the server but still it have some problem
	// it will take some [as of processing if any ongoing request]
	// or sometime it gets infinitely hang so our PORT will be locked

	// so we use timer, like if after this time shutdown didn't happen then report that, for that we use CONTEXT [a package][we pass an empty context/starting point [it is just a container to store anything]]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel timeout
	defer cancel()

	//this ctx will be passed to the server
	// err := server.Shutdown(ctx)
	// if err != nil {
	// 	slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	// }

	//SHORTFORM
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("server shutown successfully")
	// fmt.Println("server started")
	//to test fo run cmd/../main.go -config config/local.yaml
}
