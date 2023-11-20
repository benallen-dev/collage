package data

func NewSharedData() *SharedData {
	return &SharedData{
		SharedData: make(map[string]User),
	}
}
