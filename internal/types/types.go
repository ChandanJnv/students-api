package types

type Student struct {
	Id    int
	Name  string `validate:"required,alpha"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,number,min=4,max=20"`
}
