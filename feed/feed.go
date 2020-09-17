package feed

import (
	"sync"
	"time"

	"github.com/AlexKLWS/youtube-audio-stream/models"
)

type ProgressUpdateFeed struct {
	mu       sync.RWMutex
	updaters []chan models.ProgressUpdate
	closed   bool
}

func New() *ProgressUpdateFeed {
	ps := &ProgressUpdateFeed{}
	return ps
}

func (f *ProgressUpdateFeed) Subscribe() <-chan models.ProgressUpdate {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return nil
	}

	ch := make(chan models.ProgressUpdate)
	f.updaters = append(f.updaters, ch)
	return ch
}

func (f *ProgressUpdateFeed) Send(update models.ProgressUpdate) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.closed {
		return
	}

	for _, ch := range f.updaters {
		go func(ch chan models.ProgressUpdate) {
			ch <- update
		}(ch)
	}
}

func (f *ProgressUpdateFeed) Close() {
	// TODO: This is a hack to make sure all active channels send updates to clients
	// before closing
	time.Sleep(500 * time.Millisecond)

	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.closed {
		f.closed = true
		for _, ch := range f.updaters {
			go func(c chan models.ProgressUpdate) {
				// We go through all pending channel updates and "pop" the values
				// before closing channel in order to avoid panic
				for {
					noValuesPending := false
					select {
					case _, ok := <-c:
						if !ok {
							noValuesPending = true
						}
					default:
						noValuesPending = true
					}
					if noValuesPending {
						break
					}
				}
				close(c)
			}(ch)
		}
	}
}
