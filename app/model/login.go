package model

// Request body saat login
type LoginRequest struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// Response body setelah login sukses
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
