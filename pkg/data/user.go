package data

type User struct {
	Name      string
	SessionId string
	ImageUrl  string
}

func (user *User) String() string {
	return user.Name + " (" + user.SessionId + ")" + " - " + user.ImageUrl
}
