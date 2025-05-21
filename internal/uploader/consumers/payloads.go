package consumers

var (
	rawPayloads       = make(chan []byte, 256)
	processedPayloads = make(chan []byte, 256)
)
