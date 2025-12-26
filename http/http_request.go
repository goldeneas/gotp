package http

type HttpRequest struct {
	verb    string
	path    string
	headers map[string]string
	content string
}
