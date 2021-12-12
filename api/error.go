package api

type Error struct {
	Error string `json:"error" binding:"required"`
}
