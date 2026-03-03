package response

type Response[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

var (
	InvalidRequestParams = "invalid request params"
	InvalidRequestBody   = "invalid request body"
)
