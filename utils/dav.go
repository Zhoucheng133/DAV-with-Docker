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

func StartWebDAVServer(id, path, username, encryptedPassword, root string) error {
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

	server := &http.Server{
		Addr: path,
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

	ln, err := net.Listen("tcp", path)
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
	rows, err := DB.Query("SELECT id, path, username, password, root FROM config WHERE running = 1")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, path, username, password, root string
		if err := rows.Scan(&id, &path, &username, &password, &root); err == nil {
			_ = StartWebDAVServer(id, path, username, password, root)
		}
	}
}
