package bus

var (
	// leaf → project
	CheckEvents = make(chan byte, 1000)

	// project → UI
	ProjectEvents = make(chan byte, 1000)
)
