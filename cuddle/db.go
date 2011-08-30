package cuddle

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/memcache"
	"os"
)

// Rooms are stored in the datastore to be the parent entity of many Clients,
// keeping all the participants in a particular chat in the same entity group.

// Room represents a chat room.
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
	_, err := datastore.Put(c, key, client)
	if err != nil {
		return "", err
	}

	// Purge the now-invalid cache record (if it exists).
	err = memcache.Delete(c, r.Name)
	if err != nil && err != memcache.ErrCacheMiss {
		return "", err
	}

	return channel.Create(c, id)
}

func (r *Room) Send(c appengine.Context, message string) os.Error {
	var clients []Client

	_, err := memcache.JSON.Get(c, r.Name, &clients)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	if err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("Client").Ancestor(r.Key())
		_, err = q.GetAll(c, &clients)
		if err != nil {
			return err
		}
	}

	for _, client := range clients {
		err = channel.Send(c, client.ClientID, message)
		if err != nil {
			c.Errorf("sending %q: %v", message, err)
		}
	}

	return memcache.JSON.Set(c, &memcache.Item{
		Key: r.Name, Object: clients,
	})
}

// getRoom fetches a Room by name from the datastore,
// creating it if it doesn't exist already.
func getRoom(c appengine.Context, name string) (
			room *Room, err os.Error) {
	room = &Room{Name: name}
	err = datastore.RunInTransaction(c,
				func(c appengine.Context) os.Error {
		err := datastore.Get(c, room.Key(), room)
		if err == datastore.ErrNoSuchEntity {
			_, err = datastore.Put(c, room.Key(), room)
		}
		return err
	})
	return
}
