package httpserver

import (
	"net/http"
	"os"

	"github.com/valyala/fasthttp"
)

func (s *Server) faviconHandler(ctx *fasthttp.RequestCtx) {

	fileData, _ := os.ReadFile("./public/favicon.png")
	ctx.SetContentType("image/png")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(fileData)
}

func (s *Server) indexHandler(ctx *fasthttp.RequestCtx) {

	fileData, _ := os.ReadFile("./public/index.html")
	ctx.SetContentType("text/html")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(fileData)
}
