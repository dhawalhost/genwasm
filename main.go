package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dhawalhost/genwasm/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c

		log.Println("system call: %+v", oscall)
		cancel()
	}()

	if err := startServer(ctx); err != nil {
		log.Println("failed to serve:+%v\n", err)
	}
}
func startServer(ctx context.Context) (err error) {
	r := gin.Default()
	md := cors.DefaultConfig()
	md.AllowAllOrigins = true
	md.AllowHeaders = []string{"*"}
	md.AllowMethods = []string{"*"}
	r.Use(cors.New(md))
	// r.Static("/"+models.ProjectCFG.ProjectID+"/images/", handlers.UploadPath)
	// http.HandleFunc("/upload", uploadFileHandler())
	routes.InitRoutes(r)
	srv := &http.Server{
		Addr:    ":4040",
		Handler: r,
	}
	// s.ListenAndServe()
	go func() {
		// r.Run(":9000")
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()
	log.Println("server started")

	<-ctx.Done()

	log.Println("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Println("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return

}
func main() {
	run()
}
