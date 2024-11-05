package main

type UserNotFound struct{}

func (UserNotFound) Error() string {
	return "User not found"
}

type NotAuthenticated struct{}

func (NotAuthenticated) Error() string {
	return "Not Authenticated Please Login"
}
