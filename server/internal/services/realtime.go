package services

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/zenika/tilv2back/internal/repository"
	"github.com/zenika/tilv2back/internal/structures"
)

// channels in go are not meant for broadcast.
// We have to clone events from one channel to all others, which requires some plumbing.

var originalBroadcast = make(chan structures.Event)
var Broadcast = NewBroadcastServer(context.Background(), originalBroadcast)

type BroadcastServer interface {
	Subscribe() <-chan structures.Event
	Unsubscribe(<-chan structures.Event)
}

type broadcastServer struct {
	source         <-chan structures.Event
	listeners      []chan structures.Event
	addListener    chan chan structures.Event
	removeListener chan (<-chan structures.Event)
}

func (s *broadcastServer) Subscribe() <-chan structures.Event {
	newListener := make(chan structures.Event)
	s.addListener <- newListener
	return newListener
}

func (s *broadcastServer) Unsubscribe(channel <-chan structures.Event) {
	s.removeListener <- channel
}

func NewBroadcastServer(ctx context.Context, source <-chan structures.Event) BroadcastServer {
	service := &broadcastServer{
		source:         source,
		listeners:      make([]chan structures.Event, 0),
		addListener:    make(chan chan structures.Event),
		removeListener: make(chan (<-chan structures.Event)),
	}
	go service.serve(ctx)
	return service
}

func (s *broadcastServer) serve(ctx context.Context) {
	defer func() {
		for _, listener := range s.listeners {
			if listener != nil {
				close(listener)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.addListener:
			s.listeners = append(s.listeners, newListener)
		case listenerToRemove := <-s.removeListener:
			for i, ch := range s.listeners {
				if ch == listenerToRemove {
					s.listeners[i] = s.listeners[len(s.listeners)-1]
					s.listeners = s.listeners[:len(s.listeners)-1]
					close(ch)
					break
				}
			}
		case val, ok := <-s.source:
			if !ok {
				return
			}
			for _, listener := range s.listeners {
				if listener != nil {
					select {
					case listener <- val:
					case <-ctx.Done():
						return
					}

				}
			}
		}
	}
}

func SendHeartbeatToClients() {
	originalBroadcast <- structures.Event{
		Kind: "heartbeat",
		Data: structures.Post{},
	}
}

func sendEvent(id uuid.UUID, kind string) {
	article, err := repository.GetPostById(id)
	if err == nil {
		originalBroadcast <- structures.Event{
			Kind: kind,
			Data: article,
		}
	}

}
