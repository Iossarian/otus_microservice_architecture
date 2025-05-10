package billing

type Request struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

type Response struct {
	Balance float64 `json:"balance"`
}
