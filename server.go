package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

//Achtung! Only manageCurrentUrl may mutate this!
var currentUrl string = ""
var allUrls []string = make([]string, 1)

func main() {
	urlsChan := make(chan []string, 100)

	//Serve static configuration page
	//  GET config.html
	//	GET client whatnot
	//Serve configuration REST endpoint
	//	PUT /configuration/ {list of assets}
	//	Parse asset list
	//	Respond appropriately
	//Serve current display item uri
	//	GET /current/
	//	Respond w/ text/plain uri
	//	If none, 404? 500? I guess 404

	//start managing the current url
	go manageCurrentUrl(urlsChan)

	//configuration REST endpoint
	http.HandleFunc("/configure/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Recieved post request.")
		configure(req.FormValue("list"), urlsChan)
		w.Write([]byte("Success!  http://localhost:3030/config.html"))
	})

	//currentUrl REST endpoint
	http.HandleFunc("/current/", serveCurrentUrl)

	//allUrls REST endpoint
	http.HandleFunc("/all/", serveAllUrls)

	//static file server
	// configfs := http.FileServer(http.Dir("static/"))
	// http.Handle("/config/", configfs)

	//static file server
	clientfs := http.FileServer(http.Dir("client"))
	http.Handle("/", clientfs)

	//Listen for connections and serve
	log.Println("Listening...")
	err := http.ListenAndServe(":3030", nil)
	if err != nil {
		log.Println(err.Error())
	}
}

func configure(list string, urlsChan chan []string) {
	log.Println("configuring")
	resourceList := list
	allUrls = parseResourceList(resourceList)
	log.Println("configured")
}

func parseResourceList(resourceList string) []string {
	//TODO: Check for invalid URLs???????
	log.Println("Parsing URLs.")
	firstPassSlice := strings.Split(resourceList, "\n")

	finalSlice := make([]string, len(firstPassSlice))
	finalSliceIndex := 0
	for i := range firstPassSlice {
		if firstPassSlice[i] != "" {
			finalSlice[finalSliceIndex] = firstPassSlice[i]
			finalSliceIndex++
		}
	}
	log.Printf("parsed %d urls \n", finalSliceIndex)
	return finalSlice[:finalSliceIndex]
}

func manageCurrentUrl(urlsChan chan []string) {
	for {
		if len(allUrls) > 0 {
			thisUrl := allUrls[time.Time.Minute(time.Now())%len(allUrls)]
			currentUrl = thisUrl
			log.Printf("Set current url to %s\n", thisUrl)
		}
		time.Sleep(1 * time.Minute)
	}
}

func serveCurrentUrl(w http.ResponseWriter, req *http.Request) {
	if currentUrl == "" {
		w.WriteHeader(404)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(currentUrl))
	}
}

func serveAllUrls(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(makeUrlsString(allUrls)))
}

func makeUrlsString(urls []string) string {
	//we want a newline delimited list of urls
	urlsString := ""
	for i := range urls {
		urlsString = urlsString + urls[i] + "\n"
	}
	return urlsString
}
