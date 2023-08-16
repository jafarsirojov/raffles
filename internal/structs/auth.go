package structs

type AuthPayload struct {
	Token string `json:"token"`
}

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUp struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Login     string `json:"loginEmail"`
	Password  string `json:"password"`
	Token     string `json:"-"`
}
