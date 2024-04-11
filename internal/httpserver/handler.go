package httpserver

import (
	"encoding/json"
	"errors"
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
	// index app
	case "/":
		s.indexHandler(ctx)
	case "/favicon.ico":
		s.faviconHandler(ctx)

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
		// s.auth(s.messageEventsHandler)(ctx)
		s.messageEventsHandler(ctx)

	// case "/settoken":
	// 	ctx.Response.Header.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI5OTg5MjgsInVzZXJfaWQiOiIxIiwidXNlcl9uYW1lIjoiYWRtaW4ifQ.QOfJoviq81XI3N78547XMV_UfL2a9HMqikTDcRVy3zM")
	// 	respondOkJSON(ctx)

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

		claims, err := s.parseToken(string(ctx.Request.Header.Peek("token")))
		if err != nil {
			respondWithError(ctx, http.StatusBadRequest, err.Error())
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

// libs
// .
// .
func (s *Server) parseToken(tokenString string) (jwt.MapClaims, error) {

	if tokenString == "" {
		return nil, errors.New("authorization token invalid")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.tokenSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("authorization token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("authorization token invalid")
	}

	return claims, nil
}

// common lib
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
