package goods

type Good struct {
	Id       int
	Name     string
	Size     string
	UniqCode int
}

type RemainsDTO struct {
	Name             string      `json:"name"`
	Size             string      `json:"size"`
	StorageAvailable map[int]int `json:"storage_available"`
}

type ReservedDTO struct {
	UniqId         int              `json:"uniq_id"`
	Storages       []map[string]int `json:"storages"`
	AdditionalInfo string           `json:"additional_info,omitempty"`
}

type ReleasedDTO struct {
	UniqId         int    `json:"uniq_id"`
	AdditionalInfo string `json:"additional_info,omitempty"`
}
