package repository

import (
	"cryptopricealerter/internal/models"
	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(alert *models.Alert) error
	GetAll() ([]*models.Alert, error)
	GetByID(id uint) (*models.Alert, error)
	Delete(id uint) error
	MarkTriggered(id uint) error
}

type alertRepo struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepo{db: db}
}

func (r *alertRepo) Create(alert *models.Alert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepo) GetAll() ([]*models.Alert, error) {
	var alerts []*models.Alert
	err := r.db.Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *alertRepo) GetByID(id uint) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.First(&alert, id).Error
	if err != nil {
		return nil, err
	}
	return &alert, err
}

func (r *alertRepo) Delete(id uint) error {
	return r.db.Delete(&models.Alert{}, id).Error
}

func (r *alertRepo) MarkTriggered(id uint) error {
	return r.db.Model(&models.Alert{}).Where("ID = ?", id).Update("Triggered", true).Error
}