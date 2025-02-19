package main

import (
	"fmt"
	"net/http"

	"google-trends-api/internal/config"
	"google-trends-api/src/routers"
)

func main() {
	// http.HandleFunc("/", handlers.Main_handler)
	// http.HandleFunc("/kigo", handlers.handle_kigo)
	// http.HandleFunc("/", handlers.Main_handler)
	// mux := http.NewServeMux()
	// routers.Main_route_register(mux)
	config.Init_env()

	router := routers.InitRouter()

	server := &http.Server{
		Addr:         ":" + fmt.Sprint(config.APP_CONFIG.API_PORT),
		Handler:      router,
		ReadTimeout:  config.APP_CONFIG.ReadTimeout,
		WriteTimeout: config.APP_CONFIG.WriteTimeout,
	}

	go func() {

		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}
	}()

	select {}
}
