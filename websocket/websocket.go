package websocket

import (
	"github.com/kataras/iris/context"

	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gobwas"
	"github.com/kataras/neffos/gorilla"
)

var (
	// GorillaUpgrader is an upgrader type for the gorilla/websocket subprotocol implementation.
	// Should be used on `New` to construct the websocket server.
	GorillaUpgrader = gorilla.Upgrader
	// GobwasUpgrader is an upgrader type for the gobwas/ws subprotocol implementation.
	// Should be used on `New` to construct the websocket server.
	GobwasUpgrader = gobwas.Upgrader
	// DefaultGorillaUpgrader is a gorilla/websocket Upgrader with all fields set to the default values.
	DefaultGorillaUpgrader = gorilla.DefaultUpgrader
	// DefaultGobwasUpgrader is a gobwas/ws Upgrader with all fields set to the default values.
	DefaultGobwasUpgrader = gobwas.DefaultUpgrader
	// New constructs and returns a new websocket server.
	// Listens to incoming connections automatically, no further action is required from the caller.
	// The second parameter is the "connHandler", it can be
	// filled as `Namespaces`, `Events` or `WithTimeout`, same namespaces and events can be used on the client-side as well,
	// Use the `Conn#IsClient` on any event callback to determinate if it's a client-side connection or a server-side one.
	//
	// See examples for more.
	New = neffos.New

	// GorillaDialer is a `Dialer` type for the gorilla/websocket subprotocol implementation.
	// Should be used on `Dial` to create a new client/client-side connection.
	GorillaDialer = gorilla.Dialer
	// GobwasDialer is a `Dialer` type for the gobwas/ws subprotocol implementation.
	// Should be used on `Dial` to create a new client/client-side connection.
	GobwasDialer = gobwas.Dialer
	// DefaultGorillaDialer is a gorilla/websocket dialer with all fields set to the default values.
	DefaultGorillaDialer = gorilla.DefaultDialer
	// DefaultGobwasDialer is a gobwas/ws dialer with all fields set to the default values.
	DefaultGobwasDialer = gobwas.DefaultDialer
	// Dial establishes a new websocket client connection.
	// Context "ctx" is used for handshake timeout.
	// Dialer "dial" can be either `GorillaDialer` or `GobwasDialer`,
	// custom dialers can be used as well when complete the `Socket` and `Dialer` interfaces for valid client.
	// URL "url" is the endpoint of the websocket server, i.e "ws://localhost:8080/echo".
	// The last parameter, and the most important one is the "connHandler", it can be
	// filled as `Namespaces`, `Events` or `WithTimeout`, same namespaces and events can be used on the server-side as well.
	//
	// See examples for more.
	Dial = neffos.Dial

	// IsTryingToReconnect reports whether the returning "err" from the `Server#Upgrade`
	// is from a client that was trying to reconnect to the websocket server.
	//
	// Look the `Conn#WasReconnected` and `Conn#ReconnectTries` too.
	IsTryingToReconnect = neffos.IsTryingToReconnect

	// OnNamespaceConnect is the event name which its callback is fired right before namespace connect,
	// if non-nil error then the remote connection's `Conn.Connect` will fail and send that error text.
	// Connection is not ready to emit data to the namespace.
	OnNamespaceConnect = neffos.OnNamespaceConnect
	// OnNamespaceConnected is the event name which its callback is fired after namespace successfully connected.
	// Connection is ready to emit data back to the namespace.
	OnNamespaceConnected = neffos.OnNamespaceConnected
	// OnNamespaceDisconnect is the event name which its callback is fired when
	// remote namespace disconnection or local namespace disconnection is happening.
	// For server-side connections the reply matters, so if error returned then the client-side cannot disconnect yet,
	// for client-side the return value does not matter.
	OnNamespaceDisconnect = neffos.OnNamespaceDisconnect // if allowed to connect then it's allowed to disconnect as well.
	// OnRoomJoin is the event name which its callback is fired right before room join.
	OnRoomJoin = neffos.OnRoomJoin // able to check if allowed to join.
	// OnRoomJoined is the event name which its callback is fired after the connection has successfully joined to a room.
	OnRoomJoined = neffos.OnRoomJoined // able to broadcast messages to room.
	// OnRoomLeave is the event name which its callback is fired right before room leave.
	OnRoomLeave = neffos.OnRoomLeave // able to broadcast bye-bye messages to room.
	// OnRoomLeft is the event name which its callback is fired after the connection has successfully left from a room.
	OnRoomLeft = neffos.OnRoomLeft // if allowed to join to a room, then its allowed to leave from it.
	// OnAnyEvent is the event name which its callback is fired when incoming message's event is not declared to the ConnHandler(`Events` or `Namespaces`).
	OnAnyEvent = neffos.OnAnyEvent // when event no match.
	// OnNativeMessage is fired on incoming native/raw websocket messages.
	// If this event defined then an incoming message can pass the check (it's an invalid message format)
	// with just the Message's Body filled, the Event is "OnNativeMessage" and IsNative always true.
	// This event should be defined under an empty namespace in order this to work.
	OnNativeMessage = neffos.OnNativeMessage

	// IsSystemEvent reports whether the "event" is a system event,
	// OnNamespaceConnect, OnNamespaceConnected, OnNamespaceDisconnect,
	// OnRoomJoin, OnRoomJoined, OnRoomLeave and OnRoomLeft.
	IsSystemEvent = neffos.IsSystemEvent
	// Reply is a special type of custom error which sends a message back to the other side
	// with the exact same incoming Message's Namespace (and Room if specified)
	// except its body which would be the given "body".
	Reply = neffos.Reply
	// Marshal marshals the "v" value and returns a Message's Body.
	// The "v" value's serialized value can be customized by implementing a `Marshal() ([]byte, error) ` method,
	// otherwise the default one will be used instead ( see `SetDefaultMarshaler` and `SetDefaultUnmarshaler`).
	// Errors are pushed to the result, use the object's Marshal method to catch those when necessary.
	Marshal = neffos.Marshal
)

// SetDefaultMarshaler changes the default json marshaler.
// See `Marshal` package-level function and `Message.Unmarshal` method for more.
func SetDefaultMarshaler(fn func(v interface{}) ([]byte, error)) {
	neffos.DefaultMarshaler = fn
}

// SetDefaultUnmarshaler changes the default json unmarshaler.
// See `Message.Unmarshal` method and package-level `Marshal` function for more.
func SetDefaultUnmarshaler(fn func(data []byte, v interface{}) error) {
	neffos.DefaultUnmarshaler = fn
}

// Handler returns an Iris handler to be served in a route of an Iris application.
func Handler(s *neffos.Server) context.Handler {
	return func(ctx context.Context) {
		s.Upgrade(ctx.ResponseWriter(), ctx.Request(), func(socket neffos.Socket) neffos.Socket {
			return &socketWrapper{
				Socket: socket,
				ctx:    ctx,
			}
		})
	}
}

type socketWrapper struct {
	neffos.Socket
	ctx context.Context
}

// GetContext returns the Iris Context from a websocket connection.
func GetContext(c *neffos.Conn) context.Context {
	if sw, ok := c.Socket().(*socketWrapper); ok {
		return sw.ctx
	}
	return nil
}
