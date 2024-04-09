package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

func (s *Server) router(ctx *fasthttp.RequestCtx) {

	path := string(ctx.Path())
	switch path {
	case "/register":
		postJson(s.registerHandler)(ctx)
	case "/login":
		postJson(s.loginHandler)(ctx)
	case "/profile":
		postJson(s.auth(s.getUserProfile))(ctx)
	case "/edit_profile":
		postJson(s.auth(s.updateUserProfile))(ctx)
	case "/":
		s.auth(s.indexHandler)(ctx)
	default:
		// respondWithError(ctx, http.StatusNotFound, "Page not found")
		respondWithError(ctx, http.StatusBadRequest, "bad request")
	}
}

func (s *Server) auth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		tokenString := string(ctx.Request.Header.Peek("token"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return s.tokenSecretKey, nil
		})

		if err != nil || !token.Valid {
			respondWithError(ctx, http.StatusBadRequest, "authorization token invalid")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(ctx, http.StatusBadRequest, "authorization token invalid")
			return
		}

		username := claims["username"].(string)
		ctx.SetUserValue("username", username)

		h(ctx)
	}
}

// set content type json and handle only method post
func postJson(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		if !ctx.Request.Header.IsPost() {
			respondWithError(ctx, http.StatusMethodNotAllowed, "only for post request")
			return
		}

		ctx.SetContentType("Content-Type: application/json; charset=utf-8")

		h(ctx)
	}
}

// .
func respondWithError(ctx *fasthttp.RequestCtx, statusCode int, errorText string) {
	respData, _ := json.Marshal(errorText)
	ctx.SetStatusCode(statusCode)
	ctx.Write(respData)
}

// .
func respondJSON(ctx *fasthttp.RequestCtx, statusCode int, v any) {
	respData, _ := json.Marshal(v)
	ctx.SetStatusCode(statusCode)
	ctx.Write(respData)
}

// .
func respondOkJSON(ctx *fasthttp.RequestCtx) {
	respData, _ := json.Marshal("ok")
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(respData)
}
