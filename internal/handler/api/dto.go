package api

// --- GAMES ---

type CreateGameRequest struct {
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type UpdateGameRequest struct {
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}
