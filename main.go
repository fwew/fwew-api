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
	APIVersion:  "1.3.0",
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
	"simple_search_url": "ROOT/fwew-simple/{nav}",
	"search_reverse_url": "ROOT/fwew/r/{lang}/{local}",
	"list_url": "ROOT/list",
	"list_filter_url": "ROOT/list/{args}",
	"random_url": "ROOT/random/{n}",
	"random_filter_url": "ROOT/random/{n}/{args}",
	"number_to_navi_url": "ROOT/number/r/{num}",
	"navi_to_number_url": "ROOT/number/{word}",
	"lenition_url": "ROOT/lenition",
	"version_url": "ROOT/version",
	"name_single_url": "ROOT/name/single/{n}/{s}/{dialect}",
	"name_full_url": "ROOT/name/full/{ending}/{n}/{s1}/{s2}/{s3}/{dialect}",
	"name_alu_url": "ROOT/name/alu/{n}/{s}/{nm}/{am}/{dialect}"
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

	words, err := fwew.TranslateFromNavi(navi, true)
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

/* Used with /profanity to make it run faster */
func simpleSearchWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	navi := vars["nav"]

	words, err := fwew.TranslateFromNavi(navi, false)
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

func getSingleNames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err1 := strconv.Atoi(vars["n"])
	s, err2 := strconv.Atoi(vars["s"])
	dialect := vars["dialect"]
	d := 0

	if err1 != nil || err2 != nil {
		var m message
		m.Message = fmt.Sprintf("%s: %s", fwew.Text("invalidDecimalError"), vars["n"])
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	if dialect == "forest" {
		d = 1
	} else if dialect == "reef" {
		d = 2
	}

	names := fwew.SingleNames(n, d, s)
	json.NewEncoder(w).Encode(names)
}

func getFullNames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ending := vars["ending"]
	n, err1 := strconv.Atoi(vars["n"])
	s1, err2 := strconv.Atoi(vars["s1"])
	s2, err3 := strconv.Atoi(vars["s2"])
	s3, err4 := strconv.Atoi(vars["s3"])
	dialect := vars["dialect"]
	d := 0

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		var m message
		m.Message = fmt.Sprintf("%s: %s", fwew.Text("invalidDecimalError"), vars["n"])
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	if dialect == "forest" {
		d = 1
	} else if dialect == "reef" {
		d = 2
	}

	names := fwew.FullNames(ending, n, d, [3]int{s1, s2, s3}, true)
	json.NewEncoder(w).Encode(names)
}

func getNameAlu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err1 := strconv.Atoi(vars["n"])
	s, err2 := strconv.Atoi(vars["s"])
	noun_mode := vars["nm"]
	adj_mode := vars["am"]
	dialect := vars["dialect"]
	d := 0

	if err1 != nil || err2 != nil {
		var m message
		m.Message = fmt.Sprintf("%s: %s", fwew.Text("invalidDecimalError"), vars["n"])
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	if dialect == "forest" {
		d = 1
	} else if dialect == "reef" {
		d = 2
	}

	nm := 0
	if noun_mode == "normal noun" {
		nm = 1
	} else if noun_mode == "verb-er" {
		nm = 2
	}

	am := 0
	if adj_mode == "none" {
		am = 1
	} else if adj_mode == "any" {
		am = -1
	} else if adj_mode == "normal adjective" {
		am = 2
	} else if adj_mode == "genitive noun" {
		am = 3
	} else if adj_mode == "origin noun" {
		am = 4
	} else if adj_mode == "participle verb" {
		am = 5
	} else if adj_mode == "active participle verb" {
		am = 6
	} else if adj_mode == "passive participle verb" {
		am = 7
	}

	names := fwew.NameAlu(n, d, s, nm, am)
	json.NewEncoder(w).Encode(names)
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
	myRouter.HandleFunc("/api/fwew-simple/{nav}", simpleSearchWord)
	myRouter.HandleFunc("/api/list", listWords)
	myRouter.HandleFunc("/api/list/{args}", listWords)
	myRouter.HandleFunc("/api/random/{n}", getRandomWords)
	myRouter.HandleFunc("/api/random/{n}/{args}", getRandomWords)
	myRouter.HandleFunc("/api/number/r/{num}", searchNumberReverse)
	myRouter.HandleFunc("/api/number/{word}", searchNumber)
	myRouter.HandleFunc("/api/lenition", getLenitionTable)
	myRouter.HandleFunc("/api/version", getVersion)
	myRouter.HandleFunc("/api/name/single/{n}/{s}/{dialect}", getSingleNames)
	myRouter.HandleFunc("/api/name/full/{ending}/{n}/{s1}/{s2}/{s3}/{dialect}", getFullNames)
	myRouter.HandleFunc("/api/name/alu/{n}/{s}/{nm}/{am}/{dialect}", getNameAlu)

	log.Fatal(http.ListenAndServe(":"+config.Port, myRouter))
}

func main() {
	loadConfig()
	fwew.AssureDict()
	fwew.PhonemeDistros()
	handleRequests()
}
