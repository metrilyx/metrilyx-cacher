package httpwrappers

import(
	"io"
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"encoding/json"
)

type HttpResponseData struct {
	bytes []byte
}

func NewHttpResponseData(dataBytes []byte) HttpResponseData {
	return HttpResponseData{dataBytes}
}

func (h *HttpResponseData) AsJson(iface interface{}) error {
	return json.Unmarshal(h.bytes, iface)
}

func (h *HttpResponseData) AsString() string {
	return string(h.bytes)
}

type HTTPCall struct {
	client http.Client
}

func NewHTTPCall(secure bool) *HTTPCall {
	if secure {
		httpTransport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
		return &HTTPCall{client: http.Client{Transport: httpTransport},}
	} else {
		return &HTTPCall{client: http.Client{}}
	}
}

func (h *HTTPCall) getReponseBytes(r *http.Response, e error) ([]byte, error) {
	if e != nil {
		return make([]byte, 0), e
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func (h *HTTPCall) Get(url string) (HttpResponseData, error) {
	body, err := h.getReponseBytes(h.client.Get(url))
	if err != nil {
		return HttpResponseData{}, err
	}
	return HttpResponseData{body}, nil
}

func (h *HTTPCall) Post(url string, datatype string, data io.Reader) (HttpResponseData, error) {
	resp, err := h.getReponseBytes(h.client.Post(url, datatype, data))
	if err != nil {
		return HttpResponseData{}, err
	}
	return HttpResponseData{resp}, nil
}
