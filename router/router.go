package router

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
)

type Router struct {
	srv *http.Server
}

func New(port int32) (*Router) {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	return &Router{
		srv: srv,
	}
}

func (r *Router) ListenAndServe(ctx context.Context) error {
	http.HandleFunc("/debug/info", r.info)
	http.HandleFunc("/debug/pprof", pprof.Profile)

	go func() {
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen and serve: %s", err.Error())
		}
	}()

	select {
	case <-ctx.Done():
		if err := r.srv.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown: %s", err.Error())
		}
	}

	return nil
}

func (r *Router) info(w http.ResponseWriter, _ *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info := fmt.Sprintf(""+
		"Runtime OS: %s \n"+
		"Runtime ARCH: %s \n"+
		"Goroutines count: %d; \n"+
		"Allocated heap objects: %0.4f Mb \n"+
		"Total allocated memory for heap objects for the life of the program: %0.4f Mb \n"+
		"Total memory obtained from the OS: %0.4f Mb \n"+
		"The number of completed GC cycles: %d \n"+
		"",
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumGoroutine(),
		bToMb(m.HeapAlloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	)

	_, err := io.WriteString(w, info)
	if err != nil {
		log.Printf("failed to write the string, path: /info, err: %s", err.Error())
	}
}

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}
