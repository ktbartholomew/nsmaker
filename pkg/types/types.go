package types

type CreateNamespaceRequest struct {
	Namespace string `json:"namespace"`
	Username  string `json:"username"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
