package producer

var payloads = make(chan []byte, 256)
