package orms

import (
	"context"
	"errors"
	"log"

	"github.com/subashshakya/SFSS/models"
	"gorm.io/gorm"
)

var DatabaseConnection *gorm.DB

func GetUser(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	result := DatabaseConnection.WithContext(ctx).First(&user, id)
	return user, result.Error
}

func CreateUser(ctx context.Context, userData *models.User) (models.User, error) {
	result := DatabaseConnection.WithContext(ctx).Create(&userData)
	if result.Error != nil {
		log.Println("User Created Successfully. Id:", userData.Id)
		log.Println("Rows Affected: ", result.RowsAffected)
	}
	return *userData, result.Error
}

func CheckUserInDB(ctx context.Context, user *models.User) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Find(&user)
	exists := result.RowsAffected == 0
	return exists, result.Error
}

func DeleteUser(ctx context.Context, userData *models.User) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Delete(&userData)
	if result.Error != nil {
		return false, result.Error
	}
	if result.Error == nil && result.RowsAffected != 0 {
		return true, nil
	}
	return true, nil
}

func UpdateUser(ctx context.Context, userData *models.User) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Save(&userData)
	if result.Error != nil {
		return false, result.Error
	}
	if result.Error == nil && result.RowsAffected != 0 {
		return true, nil
	}
	return true, nil
}

func GetSecureFilesOfAUser(ctx context.Context, userId int) ([]models.SecureFile, error) {
	var secureFiles []models.SecureFile
	result := DatabaseConnection.WithContext(ctx).Where("UserId = ?", userId).Find(&secureFiles)
	if result.Error != nil {
		return nil, result.Error
	}
	return secureFiles, nil
}

func UpdateFile(ctx context.Context, secureFile *models.SecureFile) (models.SecureFile, error) {
	var updatedSecureFile models.SecureFile
	result := DatabaseConnection.WithContext(ctx).Save(&secureFile)
	updatedSecureFile = *secureFile
	if result.Error != nil {
		return updatedSecureFile, result.Error
	}
	return updatedSecureFile, nil
}

func CreateSecureFile(ctx context.Context, secureFile *models.SecureFile) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Save(&secureFile)
	if result.Error != nil || result.RowsAffected == 0 {
		return false, result.Error
	}
	return true, nil
}

func CheckIfSecureFileExists(ctx context.Context, secFile *models.SecureFile) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Find(&secFile)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func GetSecureFileById(ctx context.Context, id string) (secureFile *models.SecureFile, err error) {
	var secFile models.SecureFile
	result := DatabaseConnection.WithContext(ctx).Where("Id = ?", id).Find(&secFile)
	if result.Error != nil {
		return nil, result.Error
	}
	return &secFile, nil
}

func DeleteSecureFile(ctx context.Context, id string) (bool, error) {
	var secureFile models.SecureFile
	result := DatabaseConnection.WithContext(ctx).Where("Id = ?", id).Delete(&secureFile)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func CreateSuperSecret(ctx context.Context, supaSecret *models.SuperSecret) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Save(&supaSecret)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func GetSecretsOfAUser(ctx context.Context, userId uint) ([]models.SuperSecret, error) {
	var supaSecretsList []models.SuperSecret
	result := DatabaseConnection.WithContext(ctx).Where("UserId = ?", userId).Find(&supaSecretsList)
	if result.Error != nil {
		return nil, result.Error
	}
	return supaSecretsList, nil
}

func GetSecrect(ctx context.Context, secretId string) (*models.SuperSecret, error) {
	var supaSecret models.SuperSecret
	result := DatabaseConnection.WithContext(ctx).Where("Id = ?", secretId).Find(&supaSecret)
	if result.Error != nil {
		return nil, result.Error
	}
	return &supaSecret, nil
}

func UpdateSuperSecret(ctx context.Context, supaSecret *models.SuperSecret) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Save(&supaSecret)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func DeleteSuperSecret(ctx context.Context, supaSecret *models.SuperSecret) (bool, error) {
	result := DatabaseConnection.WithContext(ctx).Delete(&supaSecret)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func ShareFile(ctx context.Context, fileShare *models.FileSharing) error {
	err := DatabaseConnection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if fileShare.RecipientId == 0 {
			return errors.New("RecipientID cannot be zero")
		}
		if err := tx.First(&fileShare.Recipient, fileShare.RecipientId).Error; err != nil {
			return err
		}
		if fileShare.SenderId == 0 {
			return errors.New("SenderID cannot be zero")
		}
		if err := tx.First(&fileShare.Sender, fileShare.SenderId).Error; err != nil {
			return err
		}
		if err := tx.Save(&fileShare).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetFileSharesOfAUser(ctx context.Context, senderId uint) ([]*models.FileSharing, error) {
	var fileSharedForUser []*models.FileSharing
	result := DatabaseConnection.WithContext(ctx).Where("SenderId = ?", senderId).Find(&fileSharedForUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return fileSharedForUser, nil
}

func ShareSecret(ctx context.Context, secretShare *models.SecretSharing) error {
	err := DatabaseConnection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if secretShare.RecipientId == 0 {
			return errors.New("RecipientID cannot be zero")
		}
		if err := tx.First(&secretShare.Recipient, secretShare.RecipientId).Error; err != nil {
			return err
		}
		if secretShare.SenderId == 0 {
			return errors.New("SenderID cannot be zero")
		}
		if err := tx.First(&secretShare.Sender, secretShare.SenderId).Error; err != nil {
			return err
		}
		if err := tx.Save(&secretShare).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetSecretSharesOfAUser(ctx context.Context, senderId uint) ([]*models.FileSharing, error) {
	var secretSharedForUser []*models.FileSharing
	result := DatabaseConnection.WithContext(ctx).Where("SenderId = ?", senderId).Find(&secretSharedForUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return secretSharedForUser, nil
}

func SecretFileCount(ctx context.Context, secretFileCount *models.SecretFileCount) error {
	if secretFileCount.UserId == 0 {
		return errors.New("UserID cannot be zero")
	}
	result := DatabaseConnection.WithContext(ctx).Save(&secretFileCount)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("could not save the secret count")
	}
	return nil
}

func SecretPasswordCount(ctx context.Context, secretPassCount *models.SecretPasswordCount) error {
	if secretPassCount.UserId == 0 {
		return errors.New("UserID cannot be zero")
	}
	result := DatabaseConnection.WithContext(ctx).Save(&secretPassCount)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("could not save the secret count")
	}
	return nil
}
