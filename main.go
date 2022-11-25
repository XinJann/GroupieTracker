package main

import (
	"groupie-tracker/controller"
	"groupie-tracker/models"
	"log"
	"net/http"
	"strconv"
)

func main() {
	PORT := "8080"
	ADDR := "127.0.0.1"

	bandsData := &models.ApiData{}
	bandsData.FeedApi()
	bandsData.CreateCaches()

	coords := &models.ApiCoords{}

	staticFiles := http.FileServer(http.Dir("view/"))
	http.Handle("/view/", http.StripPrefix("/view/", staticFiles))

	http.HandleFunc("/", bandsData.RootHandler)

	http.HandleFunc("/map", func(webpage http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			request.ParseForm()
			if len(request.Form["ID"]) != 0 {
				id, _ := strconv.Atoi(request.Form["ID"][0])
				cities,name := bandsData.CitiesTab(id)
				coords.FeedApiCord(cities,name)
				controller.ServeFile(webpage, "map.html", coords)
			}
		} else {
			controller.ServeFile(webpage, "404.html", nil)
		}
	})

	log.Printf("[INFO] - Starting server on http://" + ADDR + ":" + PORT + "/")

	go bandsData.WaitThenRefreshApi()
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("[ERROR] - Server not started properly.\n" + err.Error())
	}
}
