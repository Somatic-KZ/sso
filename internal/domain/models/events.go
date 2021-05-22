package models

import (
	"github.com/JetBrainer/sso/internal/domain"

	"time"
)

func NewEvent(tdid PolymorphicID, actionType string) *Event {
	return &Event{
		TDID:       tdid,
		ActionType: actionType,
	}
}

type Event struct {
	ID         PolymorphicID `bson:"_id" json:"id"`
	TDID       PolymorphicID `bson:"tdid" json:"tdid"`
	ActionType string               `bson:"actionType" json:"actionType"`
	Verify     *Verify              `bson:"verify,omitempty" json:"verify,omitempty"`
	ExpiresAt  *time.Time           `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
}

func (e *Event) AddVerify(verify *Verify) {
	e.Verify = verify
}

func (e *Event) CleanVerify() {
	if e.Verify != nil {
		e.Verify.Send = false
		e.Verify.Tries = 0
		e.Verify.Generation = 0
		e.Verify.Expired = time.Time{}
		e.Verify.NextAttemptAt = time.Time{}
	}
}

func (e *Event) CleanVerifyWithStatus(status string) {
	if e.Verify != nil {
		e.CleanVerify()
		e.Verify.Status = status
	}
}

func (e *Event) FinishVerify() {
	if e.Verify != nil {
		e.CleanVerify()
		e.Verify.Status = domain.TokenStatusFinish
	}
}

func (e *Event) AddDeadline(deadline time.Time) {
	e.ExpiresAt = &deadline
}
