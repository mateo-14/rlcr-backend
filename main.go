package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nico-14/rlcr-backend/controllers"
	"github.com/Nico-14/rlcr-backend/db"
	"github.com/Nico-14/rlcr-backend/ds"
	"github.com/Nico-14/rlcr-backend/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type CustomRouter struct {
	*mux.Router
}

func (r CustomRouter) HandleController(prefix string, controller controllers.IController) {
	controller.Handle(prefix, r.Router)
}

func avoidSleep() {
	for {
		time.Sleep(time.Minute * 30)
		fmt.Println("Fetching avoid sleep")
		if res, err := http.Get(os.Getenv("ENDPOINT")); err == nil {
			fmt.Printf("Response StatusCode avoid sleep: %v\n", res.StatusCode)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	defer fmt.Println("Exit")

	//DB connect
	fbClient := db.New()

	//Services
	settingsService := services.NewSettingsService(fbClient.Client)
	usersService := services.NewUsersService(fbClient)
	services := &services.Services{SettSvc: settingsService, UsrSvc: usersService}

	//Connect to external resources
	ds.Connect(services)

	//Router init and config
	cr := CustomRouter{Router: mux.NewRouter()}
	cr.StrictSlash(true)
	cr.Use(mux.CORSMethodMiddleware(cr.Router))
	cr.Use(middlewareCors)

	//Handle router controllers
	cr.HandleController("api", controllers.NewSettingsController("/settings", services))
	cr.HandleController("api", controllers.NewAuthController("/auth", fbClient))
	cr.HandleController("api", controllers.NewOrdersController("/orders", services))
	cr.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go avoidSleep()
	fmt.Println("HTTP Server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, cr))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Headers, X-Requested-With")
			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
}
