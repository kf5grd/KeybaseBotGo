package main

import (
	"keybot/api"
)

func main() {
	kb = api.ChatAPI
	kb.SendUserMsg("dxb", "it works!")
}
