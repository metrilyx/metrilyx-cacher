package opentsdb

import(
    "fmt"
    "log"
	"regexp"
    "sync"
    "github.com/euforia/metrilyx-cacher/httpwrappers"
    "unicode"
)

const DEFAULT_SUGGEST_OPTS string = "max=16000000"

const ALPHABETS string    = "abcdefghijklmnopqrstuvwxyz"
const DIGITS string       = "0123456789"

var METADATA_TYPES = []string{"metrics", "tagk", "tagv"}

type MetadataResponse struct {
    httpwrappers.HttpResponseData
    MetadataType string
}

func IsValidMetaType(mdtype string) bool {
    for _, md := range METADATA_TYPES {
        if mdtype == md {
            return true
        }
    }
    return false
}

func FetchMetadataForType(url string, mdType string, outChan chan MetadataResponse) {
    ssl := false
    if url[:5] == "https" {
        ssl = true
    }
    mUrl := fmt.Sprintf("%s%s", url, mdType)
    
    nhc := httpwrappers.NewHTTPCall(ssl)
    httpResp, err := nhc.Get(mUrl)

    if err != nil {
        fmt.Println(err)
        return
    }
    outChan <- MetadataResponse{httpResp, mdType}
}

func FetchMetadata(url string) *MetadataCache {
    tsdbUrl := fmt.Sprintf("%s?%s", url, DEFAULT_SUGGEST_OPTS)

    mcache := NewMetadataCache()
        
    commChan := make(chan MetadataResponse)
    
    var wg sync.WaitGroup
    for _, metaType := range METADATA_TYPES {
        wg.Add(1)
        go func(mType string) {
            defer wg.Done()
            for _, letter := range ALPHABETS {
                query := fmt.Sprintf("%s&q=%c&type=", tsdbUrl, letter)
                FetchMetadataForType(query, mType, commChan)
                
                query = fmt.Sprintf("%s&q=%c&type=", tsdbUrl, unicode.ToUpper(letter))
                FetchMetadataForType(query, mType, commChan)
            }
        }(metaType)

        wg.Add(1)
        go func(mType string) {
            defer wg.Done()
            for _, num := range DIGITS {
                query := fmt.Sprintf("%s&q=%c&type=", tsdbUrl, num)
                FetchMetadataForType(query, mType, commChan)
            }
        }(metaType)
    }
    // collect results //
    go func(cChan chan MetadataResponse) {
        for i := range cChan {
            var data []string
            i.AsJson(&data)
            //fmt.Println(i.MetadataType, len(data))
            mcache.AddByType(i.MetadataType, data)
        }
    }(commChan)

    wg.Wait()
    close(commChan)
    log.Println("Cache collection complete!")
    return mcache
}

type MetadataCache struct {
	Metric map[string]interface{}
	TagKey map[string]interface{}
	TagValue map[string]interface{}
}

func NewMetadataCache() *MetadataCache {
	return &MetadataCache{make(map[string]interface{}),
						map[string]interface{}{},
						map[string]interface{}{}}
}

func (m *MetadataCache) AddByType(mdType string, data []string) {
    switch mdType {
        case "metrics":
            for _, d := range data {
                //if _,ok := m.Metric[d]; ok {
                //    continue
                //}
                m.Metric[d] = true
            }
        case "tagk":
            for _, d := range data {
                //if _,ok := m.TagKey[d]; ok {
                //    continue
                //}
                m.TagKey[d] = true
            }
        case "tagv":
            for _, d := range data {
                //if _,ok := m.TagValue[d]; ok {
                //    continue
                //}
                m.TagValue[d] = true
            }
        default:
            log.Println("ERROR: invalid metadata type:", mdType)
    }
}

func (m *MetadataCache) getMatches(dataset map[string]interface{}, query string) []string {
    re, _ := regexp.Compile(query)
    out := make([]string, 0)
    for k, _ := range dataset {
        match := re.MatchString(k)
        if match {
            out = append(out, k)
        }
    }
    return out
}

func (m *MetadataCache) SearchByType(mdType string, query string) []string {
    switch mdType {
        case "metrics":
            return m.getMatches(m.Metric, query)
        case "tagk":
            return m.getMatches(m.TagKey, query)
        case "tagv":
            return m.getMatches(m.TagValue, query)
        default:
            log.Println("ERROR: invalid metadata type:", mdType)   
    }
    return make([]string,0)
}