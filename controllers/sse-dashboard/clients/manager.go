package clients

import "sync"

var (
	mu      sync.Mutex
	clients = make(map[chan int]bool)
)

func AddClient(ch chan int) {
	mu.Lock()
	defer mu.Unlock()
	clients[ch] = true
}

func RemoveClient(ch chan int) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ch)
	close(ch)
}

func BroadCastMessage(msg int) {
	mu.Lock()
	defer mu.Unlock()
	for ch := range clients {
		select {
		case ch <- msg:
		default:
		}
	}
}
