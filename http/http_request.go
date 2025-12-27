package http

type HTTPRequest struct {
	headers map[string]string
	queries map[string]string
	verb    string
	path    string
	content string
}

func (r *HTTPRequest) Verb() string {
	return r.verb
}

func (r *HTTPRequest) Path() string {
	return r.path
}

func (r *HTTPRequest) Content() string {
	return r.content
}

func (r *HTTPRequest) Headers() map[string]string {
	return r.headers
}

func (r *HTTPRequest) Queries() map[string]string {
	return r.queries
}

func (h *HTTPRequest) IsConnectionClose() bool {
	return h.headers["Connection"] == "close"
}
