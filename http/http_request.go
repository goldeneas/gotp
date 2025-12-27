package http

type HttpRequest struct {
	headers map[string]string
	queries map[string]string
	verb    string
	path    string
	content string
}

func (h *HttpRequest) IsConnectionClose() bool {
	return h.headers["Connection"] == "close"
}
