package mrresponse

type (
	SuccessCreatedItemResponse struct {
		ItemID  string `json:"id"`
		Message string `json:"message,omitempty"`
	}

	SuccessModifyItemResponse struct {
		Message string `json:"message,omitempty"`
	}
)
