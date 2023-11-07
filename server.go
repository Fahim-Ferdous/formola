package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func enqueue(inp *Message) error {
	mq := &MessageQueue{
		Obj:        inp,
		AcquiredAt: &time.Time{},
		Model:      gorm.Model{},
	}
	err := db.Create(mq).Error
	if err != nil {
		return err
	}

	msgqueue <- mq
	return nil
}

//go:embed assets/favicon_io/*.ico
var assets embed.FS

func attachStaticRoutes(r *gin.Engine) {
	assetsFS := http.FS(assets)
	r.StaticFileFS("/favicon.ico", "assets/favicon_io/favicon.ico", assetsFS)
	r.StaticFileFS(
		"/android-chrome-192x192.png",
		"assets/favicon_io/android-chrome-192x192.png",
		assetsFS,
	)
	r.StaticFileFS(
		"/android-chrome-512x512.png",
		"assets/favicon_io/android-chrome-512x512.png",
		assetsFS,
	)
	r.StaticFileFS("/apple-touch-icon.png", "assets/favicon_io/apple-touch-icon.png", assetsFS)
	r.StaticFileFS("/favicon-16x16.png", "assets/favicon_io/favicon-16x16.png", assetsFS)
	r.StaticFileFS("/favicon-32x32.png", "assets/favicon_io/favicon-32x32.png", assetsFS)
	r.StaticFileFS("/site.webmanifest", "assets/favicon_io/site.webmanifest", assetsFS)
}

func attachRoutes(r *gin.Engine) {
	attachStaticRoutes(r)

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "form.gohtml", "123912391293123")
	})

	r.POST("/f/:form_id", func(ctx *gin.Context) {
		inp := &Message{}
		if err := ctx.Bind(inp); err != nil {
			// TODO: Use gin.Bind middleware.
			return // TODO:Log error.
		}

		formId := &FormID{}
		if err := ctx.BindUri(formId); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return // TODO:Log error.
		}

		inp.FormID = formId.Value
		if err := enqueue(inp); err != nil {
			if err == gorm.ErrForeignKeyViolated {
				// TODO: Error messages.
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": "Form does not exist."})
				return
			}
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return // TODO:Log error.
		}

		ctx.JSON(http.StatusAccepted, inp)
	})
}

func runServer() chan<- struct{} {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), Errors())
	engine.SetHTMLTemplate(tmpl)

	attachRoutes(engine)
	srv := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: engine,
	}
	shutdownServer := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()
	go func() {
		<-shutdownServer
		srv.Shutdown(context.Background())
	}()
	return shutdownServer
}
