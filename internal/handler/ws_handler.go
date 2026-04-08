package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/kaplenko/diplom/internal/runner"
	"github.com/kaplenko/diplom/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	submissionService *service.SubmissionService
	notifier          *runner.SubmissionNotifier
}

func NewWSHandler(submissionService *service.SubmissionService, notifier *runner.SubmissionNotifier) *WSHandler {
	return &WSHandler{
		submissionService: submissionService,
		notifier:          notifier,
	}
}

const wsTimeout = 5 * time.Minute

// SubmissionWS upgrades the connection to WebSocket and pushes the
// evaluation result as soon as it is ready. If the submission has
// already been evaluated, the result is sent immediately.
func (h *WSHandler) SubmissionWS(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("submission_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission ID"})
		return
	}

	sub, err := h.submissionService.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ws] upgrade failed for submission %d: %v", id, err)
		return
	}
	defer conn.Close()

	if sub.Status != "pending" {
		_ = conn.WriteJSON(sub)
		return
	}

	ch := h.notifier.Subscribe(id)
	defer h.notifier.Unsubscribe(id, ch)

	select {
	case result := <-ch:
		_ = conn.WriteJSON(result)
	case <-time.After(wsTimeout):
		_ = conn.WriteJSON(gin.H{"error": "timeout waiting for result"})
	}
}
