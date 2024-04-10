package httpserver

import (
	"app/internal/storage/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

//
// handlers message
//

type Messages struct {
	ID     int64  `json:"id"`
	FromID int64  `json:"from_id"`
	Text   string `json:"text"`
}

// type getMessagesRequest struct{}
type getMessagesResponse []Messages

func (s *Server) getMessagesHandler(ctx *fasthttp.RequestCtx) {
	const op = "httpserver.getMessages"
	log := s.log.With(slog.String("op", op))

	messages, err := s.storage.GetMessages()
	if err != nil {
		log.Error("%v", err)
		respondWithError(ctx, http.StatusInternalServerError, "server error")
		return
	}

	// convert to response model
	resp := make(getMessagesResponse, len(messages))
	for i, m := range messages {
		resp[i] = Messages{
			ID:     m.ID,
			FromID: m.FromID,
			Text:   m.Text,
		}
	}

	respondJSON(ctx, http.StatusOK, resp)
}

type sendMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=300"`
}

func (s *Server) sendMessageHandler(ctx *fasthttp.RequestCtx) {
	const op = "httpserver.sendMessageHandler"
	log := s.log.With(slog.String("op", op))

	reqData := ctx.PostBody()
	req := new(sendMessageRequest)

	if err := json.Unmarshal(reqData, req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "incorrect parameters")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "incorrect parameters")
		return
	}

	userID, _ := strconv.ParseInt(ctx.UserValue("user_id").(string), 10, 64)
	message := &models.Message{
		FromID: userID,
		Text:   req.Text,
	}

	_, err := s.storage.CreateMessage(message)
	if err != nil {
		log.Error("%v", err)
		respondWithError(ctx, http.StatusInternalServerError, "server error")
		return
	}

	//
	// send message to sse
	//

	respondOkJSON(ctx)
}

func (s *Server) messageEventsHandler(ctx *fasthttp.RequestCtx) {

	//
	ctx.SetContentType("text/event-stream; charset=utf8")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

	//
	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		var i int
		for {
			i++
			message := fmt.Sprintf("%d - the time is %v", i, time.Now())
			fmt.Fprintf(w, "data: Message: %s\n\n", message)
			fmt.Println(message)

			w.Flush()
			time.Sleep(5 * time.Second)
		}
	}))
}
