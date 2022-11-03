package repository

import (
	"context"
	"moonlay-test/domain"

	"gorm.io/gorm"
)

type listSQLRepository struct {
	db *gorm.DB
}

func NewListSQLRepository(db *gorm.DB) domain.ListRepository {
	return &listSQLRepository{db: db}
}

func (repository *listSQLRepository) Fetch(ctx context.Context) (lists []domain.List, err error) {
	result := repository.db.Where("parent_id IS NULL").Find(&lists)
	if result.Error != nil {
		err = result.Error
	}

	return lists, err
}

func (repository *listSQLRepository) GetById(ctx context.Context, id uint64) (list domain.List, err error) {
	result := repository.db.Preload("SubLists").Where("id = ?", id).Find(&list)

	if result.Error != nil {
		err = result.Error
	}

	if result.RowsAffected == 0 {
		err = gorm.ErrRecordNotFound
	}

	return list, err
}

func (repository *listSQLRepository) Store(ctx context.Context, list *domain.List) (err error) {
	result := repository.db.Create(&list)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *listSQLRepository) Update(ctx context.Context, list *domain.List) (err error) {
	result := repository.db.Model(&list).Updates(&list)
	if result.Error != nil {
		return result.Error
	}

	return
}

func (repository *listSQLRepository) Delete(ctx context.Context, id uint64) (err error) {
	result := repository.db.Where("id = ? AND parent_id IS NULL", id).Delete(&domain.List{})

	if result.Error != nil {
		return result.Error
	}

	return
}

func (repository *listSQLRepository) FetchSublist(ctx context.Context, parentId uint64) (lists []domain.List, err error) {
	result := repository.db.Where("parent_id = ?", parentId).Find(&lists)
	if result.Error != nil {
		err = result.Error
	}

	return lists, err
}

func (repository *listSQLRepository) GetSublistById(ctx context.Context, parentId uint64, id uint64) (list domain.List, err error) {
	result := repository.db.Where("id = ? AND parent_id = ?", id, parentId).Find(&list)

	if result.Error != nil {
		err = result.Error
	}

	if result.RowsAffected == 0 {
		err = gorm.ErrRecordNotFound
	}

	return list, err
}

func (repository *listSQLRepository) DeleteSublist(ctx context.Context, parentId uint64, id uint64) (err error) {
	result := repository.db.Where("id = ? AND parent_id = ?", id, parentId).Delete(&domain.List{})

	if result.Error != nil {
		return result.Error
	}

	return
}
