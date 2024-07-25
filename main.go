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
	APIVersion:  "1.5.3",
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
	"ROOT/": "Fwew API Index", 
	"ROOT/fwew/{nav}": "Search Word Na'vi -> Local (returns 2-Dimensional Word array)", 
	"ROOT/fwew/r/{lang}/{local}": "Search Word Local -> Na'vi (returns 2-Dimensional Word array)", 
	"ROOT/fwew-1d/{nav}": "search Word Na'vi -> Local (returns 1-Dimensional Word array)", 
	"ROOT/fwew-1d/r/{lang}/{local}": "Search Word Local -> Na'vi (returns 1-Dimensional Word array)'", 
	"ROOT/fwew-simple/{nav}": "Search Na'vi -> Local without checking affixes (returns 2-Dimensional Word array)", 
	"ROOT/homonyms": "List Na'vi Homonyms", 
	"ROOT/lenition": "Na'vi Lenition Table", 
	"ROOT/list": "List all Words (returns 1-Dimensional Word array)", 
	"ROOT/list/{args}": "List Words with attribute filtering", 
	"ROOT/list2/{c}/{args}": "List Words with attribute filtering and check-digraphs options", 
	"ROOT/multi-ipa": "List Words with multiple IPA values (alternative pronunciation)", 
	"ROOT/multiwordwords": "List Words that have two or more parts separated by a space", 
	"ROOT/name/alu/{n}/{s}/{nm}/{am}/{dialect}": "Generate title style name(s)", 
	"ROOT/name/full/{ending}/{n}/{s1}/{s2}/{s3}/{dialect}": "Generate Na'vi names in full canonical format", 
	"ROOT/name/single/{n}/{s}/{dialect}": "Generate single Na'vi names", 
	"ROOT/number/{word}": "Search a Na'vi number word to see the decimal and octal numeral forms", 
	"ROOT/number/r/{num}": "Search an integer number between 0 and 32767 to see the Na'vi word and octal numeral forms", 
	"ROOT/oddballs": "List Words that are canon but contradict Na'vi syllable rules", 
	"ROOT/phonemedistros": "Get Phoneme Distribution data", 
	"ROOT/random/{n}": "Get random Words", 
	"ROOT/random/{n}/{args}": "Get random Words with attribute filtering", 
	"ROOT/random2/{n}/{c}": "Get random Words with check-digraphs options", 
	"ROOT/random2/{n}/{c}/{args}": "Get random Words with attribute filtering and check-digraphs options", 
	"ROOT/reef/{i}": "Get Reef Na'vi syllables and IPA by Forest Na'vi IPA", 
	"ROOT/search/{lang}/{words}": "Search Na'vi <-> Local", 
	"ROOT/total-words": "Get the number of Words in the dictionary", 
	"ROOT/update": "Reload the dictionary cache", 
	"ROOT/valid/{i}": "Check if a given word string (e.g., name, loan word, etc.) follows all Na'vi syllable rules", 
	"ROOT/version": "Version information" 
}`
	endpointsJSON = strings.ReplaceAll(endpointsJSON, "ROOT", config.WebRoot)
	endpointsJSON = strings.ReplaceAll(endpointsJSON, "\n", "")
	endpointsJSON = strings.ReplaceAll(endpointsJSON, "\t", "")
	w.Write([]byte(endpointsJSON))
}

func searchWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	navi := vars["nav"]

	words, err := fwew.TranslateFromNaviHash(navi, true)
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

	words := fwew.TranslateToNaviHash(localized, languageCode)
	if len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func searchWord1d(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	navi := vars["nav"]

	words, err := fwew.TranslateFromNaviHash(navi, true)
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	oneDWords := []fwew.Word{}
	for _, a := range words {
		oneDWords = append(oneDWords, a...)
	}

	json.NewEncoder(w).Encode(oneDWords)
}

func searchWordReverse1d(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	languageCode := vars["lang"]
	localized := vars["local"]

	words := fwew.TranslateToNaviHash(localized, languageCode)
	if len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	oneDWords := []fwew.Word{}
	for _, a := range words {
		oneDWords = append(oneDWords, a...)
	}

	json.NewEncoder(w).Encode(oneDWords)
}

/* Used with /profanity to make it run faster */
func simpleSearchWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	navi := vars["nav"]

	words, err := fwew.TranslateFromNaviHash(navi, false)
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func searchBidirectional(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	languageCode := vars["lang"]
	inputWords := vars["words"]

	words, err := fwew.BidirectionalSearch(inputWords, true, languageCode)
	if err != nil || len(words) == 0 {

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

	words, err := fwew.List(args, uint8(1))
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func listWords2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	args := strings.Split(vars["args"], " ")
	c := strings.Split(vars["c"], " ")
	checkDigraphs := uint8(1)
	if c[0] == "maybe" {
		checkDigraphs = 0
	} else if c[0] == "false" {
		checkDigraphs = 2
	}

	words, err := fwew.List(args, checkDigraphs)
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
		json.NewEncoder(w).Encode(fwew.Text("invalidDecimalError"))
		return
	}

	args := strings.Split(vars["args"], " ")
	words, err := fwew.Random(n, args, uint8(1))
	if err != nil || len(words) == 0 {
		var m message
		m.Message = "no results"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func getRandomWords2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	c := strings.Split(vars["c"], " ")
	checkDigraphs := uint8(1)
	if c[0] == "maybe" {
		checkDigraphs = 0
	} else if c[0] == "false" {
		checkDigraphs = 2
	}
	if err != nil {
		json.NewEncoder(w).Encode(fwew.Text("invalidDecimalError"))
		return
	}

	args := strings.Split(vars["args"], " ")
	words, err := fwew.Random(n, args, checkDigraphs)
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

func update(w http.ResponseWriter, r *http.Request) {
	err := fwew.UpdateDict()
	if err != nil {
		var m message
		m.Message = "Update failed"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m)
		return
	} else {
		var m message
		m.Message = "Update successful"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(m)
		return
	}
}

func getSingleNames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err1 := strconv.Atoi(vars["n"])
	s, err2 := strconv.Atoi(vars["s"])
	dialect := vars["dialect"]
	d := 0

	if err1 != nil || err2 != nil {
		json.NewEncoder(w).Encode(fwew.Text("invalidDecimalError"))
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
		json.NewEncoder(w).Encode(fwew.Text("invalidDecimalError"))
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
		json.NewEncoder(w).Encode(fwew.Text("invalidDecimalError"))
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

func getPhonemeDistros(w http.ResponseWriter, r *http.Request) {
	a := fwew.GetPhonemeDistrosMap()
	json.NewEncoder(w).Encode(a)
}

func getMultiwordWords(w http.ResponseWriter, r *http.Request) {
	a := fwew.GetMultiwordWords()
	json.NewEncoder(w).Encode(a)
}

func getHomonyms(w http.ResponseWriter, r *http.Request) {
	a, _ := fwew.GetHomonyms()
	json.NewEncoder(w).Encode(a)
}

func getOddballs(w http.ResponseWriter, r *http.Request) {
	a, _ := fwew.GetOddballs()
	json.NewEncoder(w).Encode(a)
}

func getMultiIPA(w http.ResponseWriter, r *http.Request) {
	a, _ := fwew.GetMultiIPA()
	json.NewEncoder(w).Encode(a)
}

func getDictLen(w http.ResponseWriter, r *http.Request) {
	a, _ := fwew.GetDictSize()
	json.NewEncoder(w).Encode(a)
}

func getReefFromIpa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(fwew.ReefMe(vars["i"], false))
}

func getValidity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(fwew.IsValidNavi(vars["i"]))
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
	myRouter.HandleFunc("/api/fwew/{nav}", searchWord)
	myRouter.HandleFunc("/api/fwew/r/{lang}/{local}", searchWordReverse)
	myRouter.HandleFunc("/api/fwew-1d/{nav}", searchWord1d)
	myRouter.HandleFunc("/api/fwew-1d/r/{lang}/{local}", searchWordReverse1d)
	myRouter.HandleFunc("/api/fwew-simple/{nav}", simpleSearchWord)
	myRouter.HandleFunc("/api/homonyms", getHomonyms)
	myRouter.HandleFunc("/api/lenition", getLenitionTable)
	myRouter.HandleFunc("/api/list", listWords)
	myRouter.HandleFunc("/api/list/{args}", listWords)
	myRouter.HandleFunc("/api/list2/{c}/{args}", listWords2)
	myRouter.HandleFunc("/api/multi-ipa", getMultiIPA)
	myRouter.HandleFunc("/api/multiwordwords", getMultiwordWords)
	myRouter.HandleFunc("/api/name/alu/{n}/{s}/{nm}/{am}/{dialect}", getNameAlu)
	myRouter.HandleFunc("/api/name/full/{ending}/{n}/{s1}/{s2}/{s3}/{dialect}", getFullNames)
	myRouter.HandleFunc("/api/name/single/{n}/{s}/{dialect}", getSingleNames)
	myRouter.HandleFunc("/api/number/{word}", searchNumber)
	myRouter.HandleFunc("/api/number/r/{num}", searchNumberReverse)
	myRouter.HandleFunc("/api/oddballs", getOddballs)
	myRouter.HandleFunc("/api/phonemedistros", getPhonemeDistros)
	myRouter.HandleFunc("/api/random/{n}", getRandomWords)
	myRouter.HandleFunc("/api/random/{n}/{args}", getRandomWords)
	myRouter.HandleFunc("/api/random2/{n}/{c}", getRandomWords2)
	myRouter.HandleFunc("/api/random2/{n}/{c}/{args}", getRandomWords2)
	myRouter.HandleFunc("/api/reef/{i}", getReefFromIpa)
	myRouter.HandleFunc("/api/search/{lang}/{words}", searchBidirectional)
	myRouter.HandleFunc("/api/total-words", getDictLen)
	myRouter.HandleFunc("/api/update", update)
	myRouter.HandleFunc("/api/valid/{i}", getValidity)
	myRouter.HandleFunc("/api/version", getVersion)

	log.Fatal(http.ListenAndServe(":"+config.Port, myRouter))
}

func main() {
	loadConfig()
	fwew.StartEverything()
	handleRequests()
}
