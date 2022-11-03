package domain

import (
	"context"
	"time"
)

type List struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	ParentID    *uint64   `json:"parent_id"`
	Title       string    `json:"title" gorm:"size:255" validate:"required,alphanum,max=100"`
	Description string    `json:"description" gorm:"size:1000" validate:"required,max=1000"`
	File        string    `json:"file" gorm:"size:255;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SubLists    []List    `json:"sub_lists,omitempty" gorm:"foreignKey:ParentID;references:ID;"`
}

type ListRepository interface {
	Fetch(ctx context.Context) (lists []List, err error)
	GetById(ctx context.Context, id uint64) (List, error)
	Store(context.Context, *List) error
	Update(ctx context.Context, list *List) error
	Delete(ctx context.Context, id uint64) error

	FetchSublist(ctx context.Context, parentId uint64) (lists []List, err error)
	GetSublistById(ctx context.Context, parentId uint64, id uint64) (List, error)
	DeleteSublist(ctx context.Context, parentId uint64, id uint64) error
}

type ListUsecase interface {
	Fetch(ctx context.Context) (lists []List, err error)
	GetById(ctx context.Context, id uint64) (List, error)
	Store(ctx context.Context, list *List) error
	Update(ctx context.Context, list *List) error
	Delete(ctx context.Context, id uint64) error

	FetchSublist(ctx context.Context, parentId uint64) (lists []List, err error)
	GetSublistById(ctx context.Context, parentId uint64, id uint64) (List, error)
	StoreSublist(ctx context.Context, list *List) error
	DeleteSublist(ctx context.Context, parentId uint64, id uint64) error
}
