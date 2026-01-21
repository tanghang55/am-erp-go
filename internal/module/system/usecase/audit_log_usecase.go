package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type AuditLogUsecase struct {
	repo domain.AuditLogRepository
}

func NewAuditLogUsecase(repo domain.AuditLogRepository) *AuditLogUsecase {
	return &AuditLogUsecase{repo: repo}
}

func (uc *AuditLogUsecase) List(params domain.AuditLogListParams) ([]*domain.AuditLog, int64, error) {
	return uc.repo.List(params)
}

func (uc *AuditLogUsecase) RecordFromContext(c *gin.Context, payload AuditLogPayload) error {
	if c == nil {
		return errors.New("missing context")
	}

	traceID := c.GetHeader("X-Trace-Id")
	if traceID == "" {
		traceID = newTraceID()
	}

	userID, _ := c.Get(auth.UserIDKey)
	username, _ := c.Get(auth.UsernameKey)

	changes, _ := json.Marshal(map[string]any{"before": payload.Before, "after": payload.After})

	log := &domain.AuditLog{
		TraceID:    traceID,
		Module:     payload.Module,
		Action:     payload.Action,
		EntityType: payload.EntityType,
		EntityID:   payload.EntityID,
		Changes:    string(changes),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if userID != nil {
		if id, ok := userID.(uint64); ok {
			log.UserID = &id
		}
	}
	if username != nil {
		if name, ok := username.(string); ok {
			log.Username = name
		}
	}

	return uc.repo.Create(log)
}

type AuditLogPayload struct {
	Module     string
	Action     string
	EntityType string
	EntityID   string
	Before     any
	After      any
}

func newTraceID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	encoded := hex.EncodeToString(b[:])
	return encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:]
}
