package httpserver

type HTTPError struct {
	Code    int    `json:"_"`
	Message string `json:"error"`
	Detail  string `json:"error_description"`
}

func (e *HTTPError) Error() string {
	return e.Message
}
