package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func main() {
	path := "dir"

	// 从环境变量读取用户名和密码
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	port := "3000"

	if username == "" || password == "" {
		log.Fatal("Error: USERNAME and PASSWORD must be set")
	}

	fs := &webdav.Handler{
		FileSystem: webdav.Dir(path),
		LockSystem: webdav.NewMemLS(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		u, p, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if u != username || p != password {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}
		fs.ServeHTTP(w, req)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
