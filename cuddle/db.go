package cuddle

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"os"
)

type Room struct {
	Name string
}

type Client struct {
	ClientID string
}

func (r *Room) Key() *datastore.Key {
	return datastore.NewKey("Room", r.Name, 0, nil)
}

func (r *Room) AddClient(c appengine.Context, id string) (string, os.Error) {
	key := datastore.NewKey("Client", id, 0, r.Key())
	client := &Client{ClientID: id}
	if _, err := datastore.Put(c, key, client); err != nil {
		return "", err
	}
	return channel.Create(c, id)
}

func (r *Room) Send(c appengine.Context, msg string) os.Error {
	var clients []Client
	q := datastore.NewQuery("Client").Ancestor(r.Key())
	_, err := q.GetAll(c, &clients)
	if err != nil {
		return err
	}
	
	for _, client := range clients {
		err = channel.Send(c, client.ClientID, msg)
		if err != nil {
			c.Errorf("sending %q: %v", msg, err)
		}
	}

	return nil
}

// getRoom fetches a Room by name from the datastore,
// creating it if it doesn't exist already.
func getRoom(c appengine.Context, name string) (*Room, os.Error) {
	room := &Room{Name: name}

	err := datastore.RunInTransaction(c, func(c appengine.Context) os.Error {
		err := datastore.Get(c, room.Key(), room)
		if err == datastore.ErrNoSuchEntity {
			_, err = datastore.Put(c, room.Key(), room)
		}
		return err
	})

	return room, err
}
