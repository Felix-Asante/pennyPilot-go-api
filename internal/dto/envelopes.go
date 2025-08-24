package dto

import (
	"time"
)

type CreateEnvelopeDto struct {
	Name               string     `json:"name" validate:"required,min=3"`
	Account            string     `json:"account" validate:"required"`
	AutoAllocate       *bool      `json:"auto_allocate" validate:"required,boolean"`
	TargetAmount       *float64   `json:"target_amount" validate:"omitempty"` //becomes required if type=goal
	TargetedDate       *time.Time `json:"targeted_date" validate:"omitempty"`
	AllocationStrategy *string    `json:"allocation_strategy" validate:"required_if=AutoAllocate true,omitempty,oneof=fixed_amount percentage"`
	AllocationValue    *float64   `json:"allocation_value" validate:"required_if=AutoAllocate true,omitempty"`
}

type UpdateEnvelopeDto struct {
	Name               *string    `json:"name" validate:"omitempty"`
	AutoAllocate       *bool      `json:"auto_allocate" validate:"omitempty,boolean"`
	TargetAmount       *float64   `json:"target_amount" validate:"omitempty"`
	TargetedDate       *time.Time `json:"targeted_date" validate:"omitempty"`
	AllocationStrategy *string    `json:"allocation_strategy" validate:"required_if=AutoAllocate true,omitempty,oneof=fixed_amount percentage"`
	AllocationValue    *float64   `json:"allocation_value" validate:"required_if=AutoAllocate true,omitempty"`
	IsActive           *bool      `json:"is_active" validate:"omitempty,boolean"`
}
