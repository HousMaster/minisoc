package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

// .
// router
// .
func (s *Server) router(ctx *fasthttp.RequestCtx) {

	path := string(ctx.Path())
	switch path {
	case "/":
		// index app
		s.indexHandler(ctx)

		// user
	case "/register":
		post(s.registerHandler)(ctx)
	case "/login":
		post(s.loginHandler)(ctx)
	case "/profile":
		post(s.auth(s.getUserProfile))(ctx)
	case "/edit_profile":
		post(s.auth(s.updateUserProfile))(ctx)

		// message
	case "/send_message":
		post(s.auth(s.sendMessageHandler))(ctx)
	case "/get_messages":
		s.getMessagesHandler(ctx)
	case "/message_events":
		s.messageEventsHandler(ctx)

		// 404
	default:
		// respondWithError(ctx, http.StatusNotFound, "Page not found")
		respondWithError(ctx, http.StatusBadRequest, "bad request")
	}
}

// middlewares
// .
// .
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

		ctx.SetUserValue("user_id", claims["user_id"].(string))
		ctx.SetUserValue("user_name", claims["user_name"].(string))

		h(ctx)
	}
}

// .
func post(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		if !ctx.Request.Header.IsPost() {
			respondWithError(ctx, http.StatusMethodNotAllowed, "only for post request")
			return
		}

		h(ctx)
	}
}

//	libs
//
// .
// .
func respondJSON(ctx *fasthttp.RequestCtx, statusCode int, v any) {
	ctx.SetContentType("Content-Type: application/json; charset=utf-8")
	respData, _ := json.Marshal(v)
	ctx.SetStatusCode(statusCode)
	ctx.Write(respData)
}

// .
func respondOkJSON(ctx *fasthttp.RequestCtx) {
	respondJSON(ctx, http.StatusOK, "ok")
}

// .
func respondWithError(ctx *fasthttp.RequestCtx, statusCode int, errorText string) {
	respondJSON(ctx, statusCode, errorText)
}
