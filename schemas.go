package main

import (
	"time"

	"gorm.io/gorm"
)

type Required int

const (
	RequireFirstNam Required = 1
	RequireLastName Required = 1 << 1
	RequirePhone    Required = 1 << 2
	RequireEmail    Required = 1 << 3
	RequireCountry  Required = 1 << 4
	RequireCity     Required = 1 << 5
	RequireAddress  Required = 1 << 6
	RequireSubject  Required = 1 << 7
)

type FormID struct {
	Value uint64 `uri:"form_id" binding:"required"`
}

type Message struct {
	FirstName *string `form:"omitempty,first_name" binging:"omitempty,alphanumunicode,max=64" json:"first_name,omitempty" gorm:"size:64"`
	LastName  *string `form:"omitempty,last_name"  binging:"omitempty,alphanumunicode,max=64" json:"last_name,omitempty"  gorm:"size:64"`
	Phone     *string `form:"omitempty,phone"      binging:"omitempty,e164,max=16"            json:"phone,omitempty"      gorm:"size:16"`
	Email     *string `form:"omitempty,email"      binging:"omitempty,email,max=320"          json:"email,omitempty"      gorm:"size:320"`
	Country   *string `form:"omitempty,country"    binging:"omitempty,country_code,max=2"     json:"country,omitempty"    gorm:"size:2"`
	City      *string `form:"omitempty,city"       binging:"omitempty,max=32"                 json:"city,omitempty"       gorm:"size:32"`
	Address1  *string `form:"omitempty,address1"   binging:"omitempty,max=128"                json:"address1,omitempty"   gorm:"size:128"`
	Address2  *string `form:"omitempty,address2"   binging:"omitempty,max=128"                json:"address2,omitempty"   gorm:"size:128"`
	Subject   *string `form:"omitempty,subject"    binging:"omitempty,max=256"                json:"subject,omitempty"    gorm:"size:256"`

	Body string `form:"body" binding:"required,max=4096" json:"body" gorm:"not null;size:4096"`

	Form   Form   `form:"-" json:"-"`
	FormID uint64 `form:"-" json:"form_id" gorm:"not null"`
}

type Form struct {
	gorm.Model

	User   User
	UserID int

	Config int
}

type Queue[T any] struct {
	Obj        T `gorm:"embedded"`
	AcquiredAt *time.Time

	gorm.Model
}

type User struct {
	gorm.Model

	DiscordID string
	ChannelID string
}

type MessageQueue Queue[*Message]
