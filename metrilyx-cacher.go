package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/metrilyx/metrilyx-cacher/opentsdb"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

/* endpoint to tsdb name mapping */
var MetaMap = map[string]string{
	"metrics":   "metrics",
	"tagnames":  "tagk",
	"tagvalues": "tagv",
}

type TSDBDataprovider struct {
	URI            string `json:"uri"`
	Port           int    `json:"port"`
	SearchEndpoint string `json:"search_endpoint"`
}

type DataproviderConfig struct {
	Dataprovider TSDBDataprovider `json:"dataprovider"`
}

func getTsdbUrl(cfgfile string) string {
	file, e := ioutil.ReadFile(cfgfile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var dpconfig DataproviderConfig
	json.Unmarshal(file, &dpconfig)

	return fmt.Sprintf("%s:%d%s",
		dpconfig.Dataprovider.URI,
		dpconfig.Dataprovider.Port,
		dpconfig.Dataprovider.SearchEndpoint)
}

func printHelp() {
	fmt.Printf("\n Usage:\n\n")
	flag.PrintDefaults()
	fmt.Println()
}

func initFlags() (string, string, string, int) {
	var listenAddr string
	var tsdbUrl string
	var endpoint string = "/"
	var refreshInterval int
	var configFile string

	flag.StringVar(&listenAddr, "listen-addr", ":8989", "HTTP Server Port")
	flag.IntVar(&refreshInterval, "refresh-interval", 180, "Cache refresh in seconds")
	//flag.StringVar(&endpoint, "endpoint", "/", "Endpoint prefix to serve data")
	flag.StringVar(&tsdbUrl, "url", "", "Suggest URL endpoint to OpenTSDB (e.g. http://localhost:4242/api/suggest)")
	flag.StringVar(&configFile, "config", "", "Configuration file instead of CLI options")
	flag.Parse()

	if configFile != "" {
		tsdbUrl = getTsdbUrl(configFile)
	} else {
		if tsdbUrl == "" {
			fmt.Printf("OpenTSDB url required!")
			printHelp()
			os.Exit(1)
		}
	}
	return listenAddr, tsdbUrl, endpoint, refreshInterval
}

func writeHttpResponse(writer http.ResponseWriter, data []byte, respCode int) int {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(respCode)
	writer.Write(data)
	return respCode
}

func writeHttpOptionsResponse(writer http.ResponseWriter) int {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(200)
	return 200
}

func main() {

	LISTEN_ADDR, TSDB_URL, SERVE_ENDPOINT, REFRESH_INTERVAL := initFlags()
	log.Printf("Using datasource: %s\n", TSDB_URL)
	mcache := opentsdb.NewMetadataCache()
	go func() {
		log.Println("Starting initial cache collection...")
		mcache = opentsdb.FetchMetadata(TSDB_URL)
	}()

	log.Printf("Setting caching schedule: %d secs...\n", REFRESH_INTERVAL)
	ticker := time.NewTicker(time.Duration(REFRESH_INTERVAL) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				log.Println("Running cache collection...")
				mcache = opentsdb.FetchMetadata(TSDB_URL)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	http.HandleFunc(SERVE_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		metadataType := r.URL.Path[1:]
		var respCode int

		if r.Method == "OPTIONS" {
			respCode = writeHttpOptionsResponse(w)
		} else {
			_, ok := MetaMap[metadataType]
			if !ok {
				respCode = writeHttpResponse(w, []byte(`{"error": "page not found"}`), 404)
			} else {
				params := r.URL.Query()
				if val, ok := params["q"]; ok {

					results := mcache.SearchByType(MetaMap[metadataType], val[0])
					bytes, err := json.Marshal(results)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					respCode = writeHttpResponse(w, bytes, 200)
				} else {
					respCode = writeHttpResponse(w, []byte(`{"error": "'q' param required"}`), 500)
				}
			}
		}
		log.Printf("%s %d %s %s\n", r.Method, respCode, r.RequestURI, r.RemoteAddr)
	})

	log.Printf("Starting server %s%s", LISTEN_ADDR, SERVE_ENDPOINT)
	log.Fatal(http.ListenAndServe(LISTEN_ADDR, nil))
}
