package beeep

func ExampleBeep() {
	Beep(DefaultFreq, DefaultDuration)
}

func ExampleNotify() {
	var action [][]string
	Notify("Title", "MessageBody", "assets/information.png", action)
}

func ExampleAlert() {
	var action [][]string
	Alert("Title", "MessageBody", "assets/warning.png", action)
}
