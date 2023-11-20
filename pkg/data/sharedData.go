package data

import (
	"sync"
)

type SharedData struct {
	sync.Mutex
	SharedData map[string]User
}


func (sd *SharedData) GetUser(sessionId string) (User, bool) {
	sd.Lock()
	defer sd.Unlock()
	user, ok := sd.SharedData[sessionId]
	return user, ok
}

func (sd *SharedData) DeleteUser(sessionId string) {
	sd.Lock()
	defer sd.Unlock()
	delete(sd.SharedData, sessionId)
}

func (sd *SharedData) UpdateUser(user User) {
	sd.Lock()
	defer sd.Unlock()
	sd.SharedData[user.SessionId] = user
}

func (sd *SharedData) GetUsers() []User {
	sd.Lock()
	defer sd.Unlock()
	users := make([]User, 0, len(sd.SharedData))
	for _, user := range sd.SharedData {
		users = append(users, user)
	}
	return users
}
