package service

type ServiceInput struct {
	ID int `json:"id"`
}

type ServiceOutput struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}