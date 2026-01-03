package repository

import (
	"rootme-bot/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository (db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByRootMeID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("root_me_id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateOrUpdate(user *models.User) error {
	return r.db.Where(models.User{RootMeID: user.RootMeID}).Attrs(models.User{DiscordName: user.DiscordName}).FirstOrCreate(user).Error
}

func (r *UserRepository) UpdateScore(userID string, score int, lastChallenge string) error {
	return r.db.Model(&models.User{}).Where("root_me_id = ?", userID).Updates(map[string]interface{}{
			"last_score":     score,
			"last_challenge": lastChallenge,
		}).Error
}
