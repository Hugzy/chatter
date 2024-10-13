package main

type UserNotFound struct{}

func (UserNotFound) Error() string {
	return "User not found"
}
