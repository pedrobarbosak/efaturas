package models

type Values struct {
	Success bool  `json:"success"`
	Benefit int64 `json:"benefit"`
	Others  int64 `json:"others"`
}
