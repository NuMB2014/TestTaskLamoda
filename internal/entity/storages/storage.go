package storages

type Storage struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	RawAvailable string `json:"-"`
	Available    bool   `json:"available"`
}
