package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func serve(cmd *cobra.Command, args []string) error {
	addr, err := cmd.Flags().GetString("address")
	if err != nil {
		return err
	}
	configDir, err := cmd.Flags().GetString("configdir")
	if err != nil {
		return err
	}

	server := server{
		cfgFolder: configDir,
	}
	httpServer := http.Server{
		Addr:    addr,
		Handler: server,
	}
	fmt.Println("listening on " + addr)
	return httpServer.ListenAndServe()
}

type server struct {
	cfgFolder string
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := hostname(r.Host)
	h := s.findHandler(host)
	if h == nil {
		msg := "no handler found for " + host
		fmt.Println(msg)
		w.WriteHeader(502)
		w.Write([]byte(msg))
		return
	}

	h.ServeHTTP(w, r)
}

func hostname(host string) string {
	colon := strings.IndexByte(host, ':')
	if colon == -1 {
		return host
	}
	return host[0:colon]
}

func (s server) findHandler(fullhost string) http.Handler {
	host := fullhost
	for {
		if h := s.handler(host, fullhost); h != nil {
			return h
		}

		dot := strings.LastIndexByte(host, '.')
		if dot == -1 {
			return nil
		}
		host = host[0:dot]
	}
}

func (s server) handler(host, fullhost string) http.Handler {
	content, err := ioutil.ReadFile(filepath.Join(s.cfgFolder, host))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("could not read the config of %s: %v", host, err)
		}
		return nil
	}

	cfg := strings.TrimSpace(string(content))

	if _, err = strconv.Atoi(cfg); err == nil {
		return newSingleHostReverseProxy(fullhost + ":" + cfg)
	}
	return nil
}

func newSingleHostReverseProxy(host string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.Host = host
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}
