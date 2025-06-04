package utils

type User struct {
	UserID   string  `json:"user_id"`
	UserName string  `json:"user_name"`
	PhoneNo  string  `json:"phone_no"`
	Age      float64 `json:"age"`
	Password string  `json:"password"`
	Gmail    string  `json:"gmail"`
}

type UserLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type ForgetPasswordRequest struct {
	Gmail string `json:"gmail"`
	Phone string `json:"phone"`
}
