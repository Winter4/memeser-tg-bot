package main

import (
	"memeserBot/pkg/telegram"
)

func main() {

	myBot := telegram.NewBot()
	myBot.Start()
}
