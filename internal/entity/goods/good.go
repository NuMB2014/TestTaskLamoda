package goods

type Good struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	UniqCode int    `json:"uniq_code"`
}

type RemainsDTO struct {
	Name             string      `json:"name"`
	Size             string      `json:"size"`
	StorageAvailable map[int]int `json:"storage_available"`
}

type ReservedDTO struct {
	UniqCode       int              `json:"uniq_code"`
	Storages       []map[string]int `json:"storages"`
	AdditionalInfo string           `json:"additional_info,omitempty"`
}

type ReleasedDTO struct {
	UniqCode       int    `json:"uniq_code"`
	AdditionalInfo string `json:"additional_info,omitempty"`
}
