package httpserver

import (
	"app/internal/storage"
	"app/internal/storage/models"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

//
// user authorization
//

type registerRequest struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=5,max=20"`
}
type registerResponse struct {
	ID int64 `json:"id"`
}

func (s *Server) registerHandler(ctx *fasthttp.RequestCtx) {

	// const op = "httpserver.registerHandler"
	// log := s.log.With(slog.String("op", op))

	reqData := ctx.PostBody()
	req := new(registerRequest)

	if err := json.Unmarshal(reqData, req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "incorrect parameters")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "incorrect parameters")
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}
	userID, err := s.storage.CreateUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			respondWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(ctx, http.StatusInternalServerError, "server error")
		return
	}

	//
	respondJSON(ctx, http.StatusOK, registerResponse{
		ID: userID,
	})
}

//

type loginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=5,max=20"`
}
type loginResponse struct {
	ID int64 `json:"id"`
}

func (s *Server) loginHandler(ctx *fasthttp.RequestCtx) {

	const op = "httpserver.loginHandler"
	log := s.log.With(slog.String("op", op))

	reqData := ctx.PostBody()
	req := new(loginRequest)

	if err := json.Unmarshal(reqData, req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "incorrect parameters")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "incorrect parameters")
		return
	}

	user, err := s.storage.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			respondWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(ctx, http.StatusInternalServerError, "server error")
		return
	}

	if user.Password != req.Password {
		respondWithError(ctx, http.StatusBadRequest, "authorisation error")
		return
	}

	// create jwt token

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = fmt.Sprintf("%d", user.ID)
	claims["user_name"] = req.Username
	claims["exp"] = time.Now().Add(time.Hour * 60).Unix()
	tokenString, err := token.SignedString(s.tokenSecretKey)
	if err != nil {
		log.Error("failed create jwt token %w", err)
		respondWithError(ctx, http.StatusInternalServerError, "server error")
	}

	//
	ctx.Response.Header.Set("token", tokenString)

	//
	respondJSON(ctx, http.StatusOK, loginResponse{
		ID: user.ID,
	})
}

//
// user profile
//

type userProfileRequest struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
}
type userProfileResponse struct {
	Description string `json:"description" validate:"required,max=100"`
}

func (s *Server) getUserProfile(ctx *fasthttp.RequestCtx) {

	// const op = "httpserver.getUserProfile"
	// log := s.log.With(slog.String("op", op))

	reqData := ctx.PostBody()
	req := new(userProfileRequest)

	if err := json.Unmarshal(reqData, req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "incorrect parameters")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "incorrect parameters")
		return
	}

	userProfile, err := s.storage.GetUserProfile(req.Username)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			respondWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, storage.ErrUserDescriptionIsEmpty) {
			respondWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(ctx, http.StatusInternalServerError, "server error")
		return
	}

	respondJSON(ctx, http.StatusOK, userProfileResponse{
		Description: userProfile.Description,
	})
}

//

type userProfileUpdateRequest struct {
	Username    string `json:"username" validate:"required,min=5,max=20"`
	Description string `json:"description" validate:"required,max=100"`
}

func (s *Server) updateUserProfile(ctx *fasthttp.RequestCtx) {

	// const op = "httpserver.updateUserProfile"
	// log := s.log.With(slog.String("op", op))

	reqData := ctx.PostBody()
	req := new(userProfileUpdateRequest)

	if err := json.Unmarshal(reqData, req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "incorrect parameters")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "incorrect parameters")
		return
	}

	if ctx.UserValue("user_name").(string) != req.Username {
		respondWithError(ctx, http.StatusBadRequest, "access error")
		return
	}

	userProfile := models.UserProfile{
		Username:    req.Username,
		Description: req.Description,
	}
	err := s.storage.UpdateUserProfile(userProfile)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			respondWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(ctx, http.StatusInternalServerError, "server error")
		return
	}

	respondOkJSON(ctx)
}

// func (s *Server) indexHandler(ctx *fasthttp.RequestCtx) {

// 	// const op = "httpserver.indexHandler"
// 	// log := s.log.With(slog.String("op", op))

// 	username := ctx.UserValue("user_name").(string)
// 	respondJSON(ctx, http.StatusOK, fmt.Sprintf("Hello, %s!", username))
// }
