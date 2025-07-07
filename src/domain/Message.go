package domain

type Message struct {
	IdUser int    `json:"id_user"`
	State  string `json:"state"`
	Device string `json:"device"`
}
