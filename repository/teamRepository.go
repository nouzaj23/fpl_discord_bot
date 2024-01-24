package repository

import (
	"fpl_discord_bot/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Find(id uint) (*models.Team, error)
	Create(team models.Team) error
	Update(team models.Team) error
}

func NewTeamRepository(db *gorm.DB) *GormTeamRepository {
	return &GormTeamRepository{db}
}

type GormTeamRepository struct {
	db *gorm.DB
}

func (r *GormTeamRepository) Find(id uint) (*models.Team, error) {
	var team models.Team
	result := r.db.First(&team, id)
	return &team, result.Error
}

func (r *GormTeamRepository) Create(team models.Team) error {
	return r.db.Create(&team).Error
}

func (r *GormTeamRepository) Update(team models.Team) error {
	return r.db.Save(&team).Error
}
