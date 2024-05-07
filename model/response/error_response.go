package response

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}
