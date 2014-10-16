package main

import(
    "os"
    "fmt"
    "log"
    "flag"
    "net/http"
    "time"
    "encoding/json"
    "github.com/euforia/metrilyx-cacher/opentsdb"
)

/* endpoint to tsdb name mapping */
var MetaMap = map[string]string{"metrics": "metrics", "tagnames": "tagk", "tagvalues": "tagv"}

func initFlags() (string, string, string, int) {
    var listenAddr string
    var tsdbUrl string
    var endpoint string = "/"
    var refreshInterval int
    //var configFile string

    flag.StringVar(&listenAddr, "listen-addr", ":8989", "HTTP Server Port")
    flag.IntVar(&refreshInterval, "refresh-interval", 180, "Cache refresh in seconds")
    //flag.StringVar(&endpoint, "endpoint", "/", "Endpoint prefix to serve data")
    flag.StringVar(&tsdbUrl, "url", "", "Suggest URL endpoint to OpenTSDB (e.g. http://localhost:4242/api/suggest)")
    //flag.StringVar(&configFile, "-config", "metrilyx-cacher.json", "Configuration file")
    flag.Parse()
    
    if tsdbUrl == "" {
        fmt.Printf("OpenTSDB url required!")
        fmt.Printf("\n Usage:\n\n")
        flag.PrintDefaults()
        fmt.Println()
        os.Exit(1)
    }
    return listenAddr, tsdbUrl, endpoint, refreshInterval
}

func main() {

    LISTEN_ADDR, TSDB_URL, SERVE_ENDPOINT, REFRESH_INTERVAL := initFlags()

    mcache := opentsdb.NewMetadataCache()
    go func() {
        log.Println("Starting initial cache collection...")
        mcache = opentsdb.FetchMetadata(TSDB_URL)
    }()

    log.Printf("Setting caching schedule: %d...\n", REFRESH_INTERVAL)
    ticker := time.NewTicker(time.Duration(REFRESH_INTERVAL) * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
           select {
            case <- ticker.C:
                // do stuff
                log.Println("Running cache collection...")
                mcache = opentsdb.FetchMetadata(TSDB_URL)
            case <- quit:
                ticker.Stop()
                return
            }
        }
    }()

    http.HandleFunc(SERVE_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.RemoteAddr)

        metadataType := r.URL.Path[1:]
        _, ok := MetaMap[metadataType]
        if !ok {
            w.WriteHeader(404)
            return
        }
        params := r.URL.Query()
        if val,ok := params["q"]; ok {
            results := mcache.SearchByType(MetaMap[metadataType], val[0])
            bytes, err := json.Marshal(results)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type", "application/json")
            w.Write(bytes)
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(500)
            w.Write([]byte(`{"error": "'q' param required"}`))
        }
    })

    log.Printf("Starting server %s%s", LISTEN_ADDR, SERVE_ENDPOINT)
    log.Fatal(http.ListenAndServe(LISTEN_ADDR, nil)) 
}
