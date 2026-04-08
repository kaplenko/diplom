package runner

import (
	"sync"

	"github.com/kaplenko/diplom/internal/models"
)

// SubmissionNotifier delivers real-time updates to WebSocket
// subscribers when a submission evaluation completes.
type SubmissionNotifier struct {
	mu        sync.Mutex
	listeners map[int64][]chan *models.Submission
}

func NewSubmissionNotifier() *SubmissionNotifier {
	return &SubmissionNotifier{
		listeners: make(map[int64][]chan *models.Submission),
	}
}

func (n *SubmissionNotifier) Subscribe(submissionID int64) chan *models.Submission {
	n.mu.Lock()
	defer n.mu.Unlock()

	ch := make(chan *models.Submission, 1)
	n.listeners[submissionID] = append(n.listeners[submissionID], ch)
	return ch
}

func (n *SubmissionNotifier) Unsubscribe(submissionID int64, ch chan *models.Submission) {
	n.mu.Lock()
	defer n.mu.Unlock()

	listeners := n.listeners[submissionID]
	for i, l := range listeners {
		if l == ch {
			n.listeners[submissionID] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
	if len(n.listeners[submissionID]) == 0 {
		delete(n.listeners, submissionID)
	}
}

func (n *SubmissionNotifier) Publish(sub *models.Submission) {
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, ch := range n.listeners[sub.ID] {
		select {
		case ch <- sub:
		default:
		}
	}
	delete(n.listeners, sub.ID)
}
