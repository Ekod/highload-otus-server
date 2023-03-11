package domain

//User доменная сущность
type User struct {
	Id        int
	FirstName string
	LastName  string
	Age       int
	Gender    string
	Interests string
	City      string
	Email     string
	Password  string
}
