// envelope = goals
package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Envelope struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;column:id;index"`
	Name          string         `gorm:"column:name;not null"`
	AccountID     uuid.UUID      `gorm:"column:account_id;index"`
	CurrentAmount float64        `gorm:"column:current_amount;not null;default:0.00"`
	TargetAmount  float64        `gorm:"column:target_amount;not null;default:0.00"`
	AutoAllocate  bool           `gorm:"column:auto_allocate;not null;default:false"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	AllocationRule *AllocationRule `gorm:"foreignKey:TargetID;references:ID;constraint:OnDelete:CASCADE"`
}

type EnvelopeModel struct {
	db *gorm.DB
}

func NewEnvelopeModel(db *gorm.DB) *EnvelopeModel {
	return &EnvelopeModel{db: db}
}

func (em *EnvelopeModel) Create(ctx context.Context, envelope *Envelope, tx *gorm.DB) error {
	db := getTxDB(em.db, tx)

	if err := db.WithContext(ctx).Create(envelope).Error; err != nil {
		return err
	}

	return nil
}

func (em *EnvelopeModel) Save(ctx context.Context, envelope *Envelope, tx *gorm.DB) error {
	db := getTxDB(em.db, tx)

	if err := db.WithContext(ctx).Save(envelope).Error; err != nil {
		return err
	}

	return nil
}

func (em *EnvelopeModel) Delete(ctx context.Context, envelope *Envelope, tx *gorm.DB) error {
	db := getTxDB(em.db, tx)

	if err := db.WithContext(ctx).Delete(envelope).Error; err != nil {
		return err
	}

	return nil
}

func (em *EnvelopeModel) GetAllByAccountID(ctx context.Context, accountID uuid.UUID, tx *gorm.DB) ([]*Envelope, error) {
	db := getTxDB(em.db, tx)

	var envelopes []*Envelope
	if err := db.WithContext(ctx).Where("account_id = ?", accountID).Find(&envelopes).Error; err != nil {
		return nil, err
	}

	return envelopes, nil
}

func (em *EnvelopeModel) GetByID(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*Envelope, error) {
	db := getTxDB(em.db, tx)

	var envelope *Envelope
	if err := db.WithContext(ctx).Where("id = ?", id).First(&envelope).Error; err != nil {
		return nil, err
	}

	return envelope, nil
}

func (em *EnvelopeModel) GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID, tx *gorm.DB) (*Envelope, error) {
	db := getTxDB(em.db, tx)

	var envelope *Envelope
	if err := db.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&envelope).Error; err != nil {
		return nil, err
	}

	return envelope, nil
}
