package repository

import (
	"cryptopricealerter/internal/alert"

	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(alert *alert.Alert) error
	GetAll() ([]*alert.Alert, error)
	GetByID(id uint) (*alert.Alert, error)
	Delete(id uint) error
	MarkTriggered(id uint) error
}

type alertRepo struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepo{db: db}
}

func (r *alertRepo) Create(alert *alert.Alert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepo) GetAll() ([]*alert.Alert, error) {
	var alerts []*alert.Alert
	err := r.db.Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *alertRepo) GetByID(id uint) (*alert.Alert, error) {
	var alert alert.Alert
	if err := r.db.First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *alertRepo) Delete(id uint) error {
	return r.db.Delete(&alert.Alert{}, id).Error
}

func (r *alertRepo) MarkTriggered(id uint) error {
	return r.db.Model(&alert.Alert{}).Where("ID = ?", id).Update("Triggered", true).Error
}
