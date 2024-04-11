package eventserver

import (
	"context"
	"log/slog"
)

type Event struct {
	Type   string `json:"type"`
	FromID string `json:"from_id"`
	Text   string `json:"text"`
}

type EventClient struct {
	UserID     string
	RemoteAddr string
	WriteChan  chan *Event
}

type EventServer struct {
	NewEventChan         chan *Event
	NewConnectionClient  chan *EventClient
	DisconnectClientChan chan *EventClient
	ctx                  context.Context
	log                  *slog.Logger
}

func New(ctx context.Context, log *slog.Logger) *EventServer {
	return &EventServer{
		NewEventChan:         make(chan *Event, 100),
		NewConnectionClient:  make(chan *EventClient, 100),
		DisconnectClientChan: make(chan *EventClient, 100),
		ctx:                  ctx,
		log:                  log,
	}
}

func (s *EventServer) Run() {

	const op = "eventserver.Run"
	log := s.log.With(slog.String("op", op))

	clients := make(map[string]map[string]chan *Event, 10)

	for {
		select {
		// send events to all clients
		case event := <-s.NewEventChan:
			log.Debug("New event from id: %s, text: %s\n", event.FromID, event.Text)
			for _, connections := range clients {
				for _, conn := range connections {
					conn <- event
				}
			}

		// add client connection from list
		case c := <-s.NewConnectionClient:
			log.Debug("New connection id: %s, addr: %s\n", c.UserID, c.RemoteAddr)
			if _, ok := clients[c.UserID]; !ok {
				clients[c.UserID] = make(map[string]chan *Event)
			}

			clients[c.UserID][c.RemoteAddr] = c.WriteChan

		// remove client connection from list
		case c := <-s.DisconnectClientChan:
			log.Debug("Disconnection id: %s, addr: %s\n", c.UserID, c.RemoteAddr)
			if _, ok := clients[c.UserID]; !ok {
				continue
			}

			if _, ok := clients[c.UserID][c.RemoteAddr]; !ok {
				continue
			}

			delete(clients[c.UserID], c.RemoteAddr)

			if len(clients[c.UserID]) == 0 {
				delete(clients, c.UserID)
			}

		case <-s.ctx.Done():
			return
		}

	}
}
