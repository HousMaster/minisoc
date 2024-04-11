package httpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"app/internal/httpserver/eventserver"
	"app/internal/storage/models"

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
	s.eventserver.NewEventChan <- &eventserver.Event{
		Type:   "message",
		FromID: ctx.UserValue("user_id").(string),
		Text:   req.Text,
	}

	respondOkJSON(ctx)
}

func (s *Server) messageEventsHandler(ctx *fasthttp.RequestCtx) {

	const op = "httpserver.messageEventsHandler"
	log := s.log.With(slog.String("op", op))

	//
	// parse token
	tokenParams, err := s.parseToken(string(ctx.QueryArgs().Peek("token")))
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	//
	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {

		client := &eventserver.EventClient{
			UserID:     tokenParams["user_id"].(string),
			RemoteAddr: ctx.RemoteAddr().String(),
			WriteChan:  make(chan *eventserver.Event),
		}

		s.eventserver.NewConnectionClient <- client
		s.eventserver.NewEventChan <- &eventserver.Event{
			Type:   "connect",
			FromID: client.UserID,
			Text:   fmt.Sprintf("new connection %s:%s", client.UserID, client.RemoteAddr),
		}

		ticker := time.NewTicker(2 * time.Second)

		defer func() {
			s.eventserver.DisconnectClientChan <- client
			s.eventserver.NewEventChan <- &eventserver.Event{
				Type:   "disconnect",
				FromID: client.UserID,
				Text:   fmt.Sprintf("disconnection %s:%s", client.UserID, client.RemoteAddr),
			}
			ticker.Stop()
		}()

	LOOP:
		for {
			select {

			// ping client
			case <-ticker.C:

				fmt.Fprint(w, "data: ping\n\n")
				if err := w.Flush(); err != nil {
					log.Error("%v", err)
					break LOOP
				}

			// send event client
			case event := <-client.WriteChan:

				data, err := json.Marshal(event)
				if err != nil {
					log.Error("%v", err)
					continue
				}

				fmt.Fprintf(w, "data:%s\n\n", data)
				if err := w.Flush(); err != nil {
					log.Error("%v", err)
					break LOOP
				}

			// server shutdown (-)
			case <-ctx.Done():
				fmt.Println("ctx.Done")
				break LOOP
			}
		}
	})
}
