package main

import (
	"keybot/api"
)

func main() {
	u := api.User{Name: "dxb"}
	u.SendMessage("It works!!!")
}
