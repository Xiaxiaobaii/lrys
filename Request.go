package lrys

import (
	tool "github.com/Xiaxiaobaii/autotool"
	"bytes"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func HttpRequest(method, Url, Proxy string, proxyBool bool, header map[string]string) (data *http.Response, err error) {
	c := &http.Client{}
	if proxyBool {
		c.Transport = &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(Proxy)
			},
		}
	}
	req, _ := http.NewRequest(method, Url, bytes.NewReader([]byte{}))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	return resp, err
}

func HttpRequestSocksProxy(method, Url, Proxy string, header map[string]string) (data *http.Response, err error) {
	dialer, err := proxy.SOCKS5("tcp", Proxy, nil, proxy.Direct)
	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	c := &http.Client{Transport: httpTransport}

	req, _ := http.NewRequest(method, Url, bytes.NewReader([]byte{}))

	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	return resp, err
}

func HttpRequestVpnProxy(method, Url, Proxy string, header map[string]string) (data *http.Response, err error) {
	httpTransport := &http.Transport{}
	httpTransport.Dial = func(network, addr string) (net.Conn, error) {
		port := rand.IntN(40000) + 10001
		lAddr, err := net.ResolveTCPAddr(network, Proxy+":"+tool.Itoa(port))
		if err != nil {
			return nil, err
		}
		rAddr, err := net.ResolveTCPAddr(network, addr)
		if err != nil {
			return nil, err
		}
		conn, err := net.DialTCP(network, lAddr, rAddr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	c := &http.Client{Transport: httpTransport}

	req, _ := http.NewRequest(method, Url, bytes.NewReader([]byte{}))

	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	return resp, err
}

func HttpJsonGet(url string) (map[string]interface{}, error) {
	resp, err := HttpRequest("GET", url, "", false, make(map[string]string))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ReJson := make(map[string]interface{})
	json.Unmarshal(body, &ReJson)
	return ReJson, err
}
