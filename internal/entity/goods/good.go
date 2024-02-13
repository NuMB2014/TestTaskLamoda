package goods

type RemainsDTO struct {
	Name             string      `json:"name"`
	Size             string      `json:"size"`
	StorageAvailable map[int]int `json:"storage_available"`
}
