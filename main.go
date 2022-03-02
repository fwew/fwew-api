package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	fwew "github.com/fwew/fwew-lib/v5"
	"github.com/gorilla/mux"
)

// global instance of Config
var config Config

// global configured instance of Version
var version = Version{
	APIVersion:  "1.2.1",
	FwewVersion: fmt.Sprintf("%d.%d.%d", fwew.Version.Major, fwew.Version.Minor, fwew.Version.Patch),
	DictBuild:   fwew.Version.DictBuild,
}

// Config contains variables to be configured in the config.json file
type Config struct {
	Port    string `json:"Port"`
	WebRoot string `json:"WebRoot"`
}

// Version contains the API and Fwew version information.
type Version struct {
	APIVersion  string `json:"APIVersion"`
	FwewVersion string `json:"FwewVersion"`
	DictBuild   string `json:"DictVersion"`
}

// number represents a Na'vi number.
type number struct {
	Name    string `json:"name"`
	Octal   string `json:"octal"`
	Decimal string `json:"decimal"`
}

// message represents an error message.
type message struct {
	Message string `json:"message"`
}

// load data from config.json into Config struct
func loadConfig() {
	configFile, _ := os.Open("config.json")
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err := decoder.Decode(&config)
	if err != nil {
		// fall back to default values
		config.Port = "8080"
		config.WebRoot = "https://localhost"
	}
}

func getEndpoints(w http.ResponseWriter, r *http.Request) {
	var endpointsJSON = `{
	"search_url": "ROOT/fwew/{nav}",
	"search_reverse_url": "ROOT/fwew/r/{lang}/{local}",
	"list_url": "ROOT/list",
	"list_filter_url": "ROOT/list/{args}",
	"random_url": "ROOT/random/{n}",
	"random_filter_url": "ROOT/random/{n}/{args}",
	"number_to_navi_url": "ROOT/number/r/{num}",
	"navi_to_number_url": "ROOT/number/{word}",
	"lenition_url": "ROOT/lenition",
	"version_url": "ROOT/version"
}`
	endpointsJSON = strings.ReplaceAll(endpointsJSON, "ROOT", config.WebRoot)
	endpointsJSON = strings.ReplaceAll(endpointsJSON, " ", "")
	endpointsJSON = strings.ReplaceAll(endpointsJSON, "\n", "")
	endpointsJSON = strings.ReplaceAll(endpointsJSON, "\t", "")
	w.Write([]byte(endpointsJSON))
}

func searchWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	navi := vars["nav"]

	words, err := fwew.TranslateFromNavi(navi)
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func searchWordReverse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	languageCode := vars["lang"]
	localized := vars["local"]

	words := fwew.TranslateToNavi(localized, languageCode)
	if len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func listWords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	args := strings.Split(vars["args"], " ")

	words, err := fwew.List(args)
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func getRandomWords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		var m message
		m.Message = fmt.Sprintf("%s: %s", fwew.Text("invalidDecimalError"), vars["n"])
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	args := strings.Split(vars["args"], " ")

	words, err := fwew.Random(n, args)
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func searchNumber(w http.ResponseWriter, r *http.Request) {
	var n number
	vars := mux.Vars(r)
	d, err := fwew.NaviToNumber(vars["word"])
	if err != nil {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}
	n.Name = vars["word"]
	n.Decimal = fmt.Sprintf("%d", d)
	n.Octal = fmt.Sprintf("%#o", d)

	json.NewEncoder(w).Encode(n)
}

func searchNumberReverse(w http.ResponseWriter, r *http.Request) {
	var n number
	vars := mux.Vars(r)
	num, err := strconv.ParseInt(vars["num"], 0, 0)
	if err != nil {
		var m message
		m.Message = fmt.Sprintf("%s: %s", fwew.Text("invalidIntError"), vars["n"])
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}
	word, err := fwew.NumberToNavi(int(num))
	if err != nil {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}
	n.Name = word
	n.Decimal = fmt.Sprintf("%d", num)
	n.Octal = fmt.Sprintf("%#o", num)

	json.NewEncoder(w).Encode(n)
}

func getLenitionTable(w http.ResponseWriter, r *http.Request) {
	lenitionTableJSON := `{"kx":"k","px":"p","tx":"t","k":"h","p":"f","t":"s","ts":"s","'":"(disappears, except before ll or rr)"}`
	w.Write([]byte(lenitionTableJSON))
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(version)
}

// set the Header Content-Type to "application/json" for all endpoints
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(contentTypeMiddleware)

	myRouter.HandleFunc("/api/", getEndpoints)
	myRouter.HandleFunc("/api/fwew/r/{lang}/{local}", searchWordReverse)
	myRouter.HandleFunc("/api/fwew/{nav}", searchWord)
	myRouter.HandleFunc("/api/list", listWords)
	myRouter.HandleFunc("/api/list/{args}", listWords)
	myRouter.HandleFunc("/api/random/{n}", getRandomWords)
	myRouter.HandleFunc("/api/random/{n}/{args}", getRandomWords)
	myRouter.HandleFunc("/api/number/r/{num}", searchNumberReverse)
	myRouter.HandleFunc("/api/number/{word}", searchNumber)
	myRouter.HandleFunc("/api/lenition", getLenitionTable)
	myRouter.HandleFunc("/api/version", getVersion)

	log.Fatal(http.ListenAndServe(":"+config.Port, myRouter))
}

func main() {
	loadConfig()
	fwew.AssureDict()
	handleRequests()
}
