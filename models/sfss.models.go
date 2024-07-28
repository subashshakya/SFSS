package main

import (
	"database/sql"
	"time"
)

type User struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	FirstName   string `gorm:"not null"`
	MiddleName  sql.NullString
	LastName    string `gorm:"not null"`
	Email       string `gorm:"not null;unique"`
	Password    string `gorm:"not null"`
	PhoneNumber string `gorm:"not null"`
}

type SecureFile struct {
	Id         string    `gorm:"primaryKey"`
	FileName   string    `gorm:"not null"`
	FileData   []byte    `gorm:"not null"`
	OriginalId uint      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:current_timestamp"`
	UserId     int       `gorm:"not null"`
	User       User      `gorm:"foreignKey:UserId;references:Id"`
}

type SuperSecret struct {
	Id        string `gorm:"primaryKey"`
	Secret    string `gorm:"not null"`
	CreatedAt time.Time
	UserId    uint `gorm:"not null"`
	User      User `gorm:"foreignKey:UserId;references:Id"`
}

type FileSharing struct {
	Id          uint       `gorm:"primaryKey"`
	FileId      string     `gorm:"not null"`
	SenderId    uint       `gorm:"not null"`
	RecipientId uint       `gorm:"not null"`
	SharedAt    time.Time  `gorm:"default:current_timestamp"`
	File        SecureFile `gorm:"foreignKey:FileId;references:Id"`
	Sender      User       `gorm:"foreignKey:SenderId;references:Id"`
	Recipient   User       `gorm:"foreignKey:RecipientId;references:Id"`
}

type SecretSharing struct {
	Id          uint        `gorm:"primaryKey"`
	SecretId    string      `gorm:"not null"`
	SenderId    uint        `gorm:"not null"`
	RecipientId uint        `gorm:"not null"`
	SharedAt    time.Time   `gorm:"default:current_timestamp"`
	Secret      SuperSecret `gorm:"foreignKey:SecretId;references:Id"`
	Sender      User        `gorm:"foreignKey:SenderId;references:Id"`
	Recipient   User        `gorm:"foreignKey:RecipientId;references:Id"`
}

type SecretFileCount struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint `gorm:"not null"`
}

type SecretPasswordCount struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint `gorm:"not null"`
}
