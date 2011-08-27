package cuddle

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"os"
)

// Room represents a chat room.
// These are stored in the datastore to be the parent entity of many Clients,
// keeping all the participants in a particular chat in the same entity group.
type Room struct {
	Name string
}

func (r *Room) Key() *datastore.Key {
	return datastore.NewKey("Room", r.Name, 0, nil)
}

// Client is a participant in a chat Room.
type Client struct {
	ClientID string // the channel Client ID
}

// AddClient puts a Client record to the datastore with the Room as its
// parent, creates a channel and returns the channel token.
func (r *Room) AddClient(c appengine.Context, id string) (string, os.Error) {
	key := datastore.NewKey("Client", id, 0, r.Key())
	client := &Client{ClientID: id}
	if _, err := datastore.Put(c, key, client); err != nil {
		return "", err
	}
	return channel.Create(c, id)
}

// Send sends a message to all Clients in a Room.
func (r *Room) Send(c appengine.Context, message string) os.Error {
	var clients []Client
	q := datastore.NewQuery("Client").Ancestor(r.Key())
	_, err := q.GetAll(c, &clients)
	if err != nil {
		return err
	}
	
	for _, client := range clients {
		err = channel.Send(c, client.ClientID, message)
		if err != nil {
			c.Errorf("sending %q: %v", message, err)
		}
	}

	return nil
}

// getRoom fetches a Room by name from the datastore,
// creating it if it doesn't exist already.
func getRoom(c appengine.Context, name string) (room *Room, err os.Error) {
	room = &Room{Name: name}
	err = datastore.RunInTransaction(c, func(c appengine.Context) os.Error {
		err := datastore.Get(c, room.Key(), room)
		if err == datastore.ErrNoSuchEntity {
			_, err = datastore.Put(c, room.Key(), room)
		}
		return err
	})
	return
}
