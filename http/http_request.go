package http

type HttpRequest struct {
	verb    string
	path    string
	headers map[string]string
	content string
}

func (h *HttpRequest) IsConnectionClose() bool {
	return h.headers["Connection"] == "close"
}
