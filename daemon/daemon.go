package daemon

func Launch() {
	go NewYoutube().Start()
}
