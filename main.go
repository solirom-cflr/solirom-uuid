package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var stringRegexp = regexp.MustCompile(`([6]{3})`)
var fileMutex sync.Mutex

func main() {
	port := ":7010"

	log.Println("Listening on " + port + "...")
	log.Fatal(http.ListenAndServe(port, nil))
}

func init() {
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.HandleFunc("/{subdomain}", subdomainHandler).Methods("GET")

	http.Handle("/", mainRouter)
}

func subdomainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	subdomain := muxVariableLookup(request, "subdomain")

	newValue := uuid(subdomain)
	log.Printf("newValue = %s", newValue)
	responseWriter.Write([]byte("A" + newValue))
}

func uuid(subdomain string) (newValueStr string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	newValueStr = "1000000"

	filename := subdomain + ".txt"
	if _, err := os.Stat(filename); err == nil {
		currentValue, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error reading")
		}
		currentValueStr := strings.TrimSpace(string(currentValue))
		newValueStr = IncrementValue(currentValueStr)
	}

	err := ioutil.WriteFile(filename, []byte(newValueStr), 0644)
	if err != nil {
		log.Printf("Error writing")
	}

	return newValueStr
}

func IncrementValue(currentValueStr string) (newValueStr string) {
	newValueInt, _ := strconv.Atoi(currentValueStr)
	newValueInt++

	newValueStr = strconv.Itoa(newValueInt)

	uuidTokens := stringRegexp.FindAllStringSubmatch(newValueStr, -1)
	if uuidTokens != nil {
		newValueStr = IncrementValue(newValueStr)
	}

	return newValueStr
}

func muxVariableLookup(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
