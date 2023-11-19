package data

import (
	"sync"
)

type User struct {
	Name      string
	SessionId string
	FileName  string
}

type SharedData struct {
	sync.Mutex
	SharedData map[string]User
}

func NewSharedData() *SharedData {
	return &SharedData{
		SharedData: make(map[string]User),
	}
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
