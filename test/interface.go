package main

import "fmt"

type UserNotifier interface {
	SendMessage(user *User, message string)
}

type User struct {
	name string
	email string
	notify UserNotifier
}

type email struct {
}

type sms struct {
}

func (e *email) SendMessage(user *User, message string) {
	fmt.Println(user.name, user.email, message)
}

func (s *sms) SendMessage(user *User, message string) {
	fmt.Println(user.name, user.email, message)
}

func (user *User) Notify(message string) {
	user.notify.SendMessage(user, message)
}

func main() {
	user1 := &User{"lisa","abc@qq.com", &email{}}
	user2 := &User{"alice","qwe@qq.com", &sms{}}

	user1.Notify("hello lisa!")
	user2.Notify("hello alice!")
}
