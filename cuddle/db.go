package cuddle

import (
	"appengine"
	"appengine/datastore"
	"os"
)

type Room struct {
	Name string
}

// getRoom fetches a Room by name from the datastore,
// creating it if it doesn't exist already.
func getRoom(c appengine.Context, name string) (*Room, os.Error) {
	room := new(Room)
	key := datastore.NewKey("Room", name, 0, nil)

	err := datastore.RunInTransaction(c, func(c appengine.Context) os.Error {
		err := datastore.Get(c, key, room)
		if err == datastore.ErrNoSuchEntity {
			room.Name = name
			_, err = datastore.Put(c, key, room)
		}
		return err
	})

	return room, err
}
