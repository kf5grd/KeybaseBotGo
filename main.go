package main

import (
	"keybot/api"
)

func main() {
	u := api.Channel{Name: "dxb"}
	u.SendMessage("It works!!!")
}
