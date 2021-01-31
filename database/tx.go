package database

type Account string


// Tx: Luu thong tin transaction
type Tx struct {
	From Account `json:"from"`  // Thong tin nguoi gui
	To Account `json:"to"`      // Thong tin nguoi nhan
	Value uint `json:"value"`   // Gia tri giao dich
	Data string `json:"data"`   //
}

// IsReward: Kiem tra xem giao dich co phai la bonus
func (t Tx) IsReward() bool {
	return t.Data == "Reward"
}




