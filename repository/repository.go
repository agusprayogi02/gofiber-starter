package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(tx *gorm.DB, entity *T) error {
	return tx.Create(entity).Error
}

func (r *Repository[T]) Update(tx *gorm.DB, entity *T) error {
	return tx.Save(entity).Error
}

func (r *Repository[T]) Delete(tx *gorm.DB, entity *T) error {
	return tx.Delete(entity).Error
}

func (r *Repository[T]) CountById(tx *gorm.DB, id any) (int64, error) {
	var total int64
	err := tx.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(tx *gorm.DB, entity *T, id any) error {
	return tx.Where("id = ?", id).Take(entity).Error
}
