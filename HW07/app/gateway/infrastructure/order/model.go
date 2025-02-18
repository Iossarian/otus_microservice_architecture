package order

type Request struct {
	UserID int     `json:"user_id"`
	Price  float64 `json:"price"`
}
