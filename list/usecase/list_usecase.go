package usecase

import (
	"context"
	"moonlay-test/domain"
	"reflect"
	"time"
)

type listUsecase struct {
	ListRepository domain.ListRepository
	ContextTimeout time.Duration
}

func NewListUsecase(li domain.ListRepository, timeout time.Duration) domain.ListUsecase {
	return &listUsecase{
		ListRepository: li,
		ContextTimeout: timeout,
	}
}

func (usecase *listUsecase) Fetch(ctx context.Context) (lists []domain.List, err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	lists, err = usecase.ListRepository.Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return
}

func (usecase *listUsecase) GetById(ctx context.Context, id uint64) (list domain.List, err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	list, err = usecase.ListRepository.GetById(ctx, id)

	return
}

func (usecase *listUsecase) Store(ctx context.Context, list *domain.List) (err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	err = usecase.ListRepository.Store(ctx, list)

	return
}

func (usecase *listUsecase) Update(ctx context.Context, list *domain.List) (err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	listExisted, err := usecase.ListRepository.GetById(ctx, list.ID)
	if err != nil {
		return
	}

	if reflect.DeepEqual(listExisted, domain.List{}) {
		return domain.ErrConflict
	}

	err = usecase.ListRepository.Update(ctx, list)

	return
}

func (usecase *listUsecase) Delete(ctx context.Context, id uint64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	err = usecase.ListRepository.Delete(ctx, id)

	return
}

func (usecase *listUsecase) FetchSublist(ctx context.Context, parentId uint64) (lists []domain.List, err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	lists, err = usecase.ListRepository.FetchSublist(ctx, parentId)

	if err != nil {
		return nil, err
	}

	return
}

func (usecase *listUsecase) GetSublistById(ctx context.Context, parentId uint64, id uint64) (list domain.List, err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	list, err = usecase.ListRepository.GetSublistById(ctx, parentId, id)

	return
}

func (usecase *listUsecase) StoreSublist(ctx context.Context, list *domain.List) (err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	err = usecase.ListRepository.Store(ctx, list)

	return
}

func (usecase *listUsecase) DeleteSublist(ctx context.Context, parentId uint64, id uint64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, usecase.ContextTimeout)
	defer cancel()

	err = usecase.ListRepository.DeleteSublist(ctx, parentId, id)

	return
}
