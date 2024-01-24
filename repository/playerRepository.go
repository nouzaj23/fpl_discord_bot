package repository

import (
	"fpl_discord_bot/models"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	Find(id uint) (*models.Player, error)
	FindByName(name string) ([]models.Player, error)
	Create(player models.Player) error
	Update(player models.Player) error
}

func NewPlayerRepository(db *gorm.DB) *GormPlayerRepository {
	return &GormPlayerRepository{db}
}

type GormPlayerRepository struct {
	db *gorm.DB
}

func (r *GormPlayerRepository) Find(id uint) (*models.Player, error) {
	var player models.Player
	result := r.db.First(&player, id)
	return &player, result.Error
}

func (r *GormPlayerRepository) FindByName(name string) ([]models.Player, error) {
	var players []models.Player
	result := r.db.Where("Name LIKE ?", "%"+name+"%").Find(&players)
	return players, result.Error
}

func (r *GormPlayerRepository) Create(player models.Player) error {
	return r.db.Create(&player).Error
}

func (r *GormPlayerRepository) Update(player models.Player) error {
	return r.db.Save(&player).Error
}
