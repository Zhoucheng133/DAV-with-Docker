package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/webdav"
)

var (
	ServersMutex sync.Mutex
	Servers      = make(map[string]*http.Server)
)

func StartWebDAVServer(id, port, username, encryptedPassword, root string) error {
	ServersMutex.Lock()
	defer ServersMutex.Unlock()

	if _, exists := Servers[id]; exists {
		return nil
	}

	password, err := Decrypt(encryptedPassword, GetAESKey())
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	fs := &webdav.Handler{
		FileSystem: webdav.Dir(root),
		LockSystem: webdav.NewMemLS(),
	}

	listenAddr := fmt.Sprintf(":%s", port)
	if len(port) > 0 && port[0] == ':' {
		listenAddr = port
	}

	server := &http.Server{
		Addr: listenAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
		}),
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	Servers[id] = server

	go func() {
		_ = server.Serve(ln)
	}()

	return nil
}

func StopWebDAVServer(id string) error {
	ServersMutex.Lock()
	defer ServersMutex.Unlock()

	server, exists := Servers[id]
	if !exists {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	delete(Servers, id)
	return err
}

func InitRunningWebDAVServers() {
	rows, err := DB.Query("SELECT id, port, username, password, root FROM config WHERE running = 1")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, port, username, password, root string
		if err := rows.Scan(&id, &port, &username, &password, &root); err == nil {
			_ = StartWebDAVServer(id, port, username, password, root)
		}
	}
}
