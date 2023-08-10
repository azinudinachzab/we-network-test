// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	ID uint64
	FullName string
	PhoneNumber string
	Password string
}

type UserRecord struct {
	UserID     uint64
	LoginCount uint64
}
