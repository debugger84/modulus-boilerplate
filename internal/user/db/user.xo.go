package db

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
)

// User represents a row from 'user.user'.
type User struct {
	ID           uuid.UUID      `gorm:"column:id" json:"id"`                         //
	Name         string         `gorm:"column:name" json:"name"`                     //
	Email        string         `gorm:"column:email" json:"email"`                   //
	RegisteredAt time.Time      `gorm:"column:registered_at" json:"registeredAt"`    //
	Settings     *Settings      `gorm:"column:settings" json:"settings"`             //
	Contacts     pq.StringArray `gorm:"column:contacts;type:text[]" json:"contacts"` //
}

func (User) TableName() string {
	return UserTable
}

type UserFinder struct {
	db *gorm.DB
}

func NewUserFinder(gormDb *gorm.DB) *UserFinder {
	return &UserFinder{db: gormDb}
}

func (f *UserFinder) OneByQuery(query QueryBuilder) *User {
	var user *User
	query.Build().Limit(1).Scan(&user)

	return user
}

func (f *UserFinder) ListByQuery(query QueryBuilder, count int) []*User {
	var users []*User
	query.Build().Limit(count).Scan(&users)

	return users
}

func (f *UserFinder) CreateQuery(ctx context.Context) *UserQuery {
	return NewUserQuery(ctx, f.db)
}

type UserSaver struct {
	db *gorm.DB
}

func NewUserSaver(db *gorm.DB) *UserSaver {
	return &UserSaver{db: db}
}

func (f *UserSaver) Create(ctx context.Context, entity User) error {
	result := f.db.Table(UserTable).WithContext(ctx).Create(&entity)

	return result.Error
}

func (f *UserSaver) Update(ctx context.Context, entity User) error {
	result := f.db.Table(UserTable).WithContext(ctx).Save(&entity)

	return result.Error
}

const UserTable = `"user"."user"`

type UserQuery struct {
	db *gorm.DB
}

func NewUserQuery(ctx context.Context, db *gorm.DB) *UserQuery {
	localCopy := db.Table(UserTable).WithContext(ctx)
	query := &UserQuery{
		db: localCopy,
	}
	return query
}

func (p *UserQuery) ID(ID uuid.UUID) *UserQuery {
	p.db = p.db.Where(UserTable+".id = ?", ID)

	return p
}

func (p *UserQuery) IDIn(items []uuid.UUID) *UserQuery {
	p.db = p.db.Where(UserTable+".id IN (?)", items)

	return p
}

func (p *UserQuery) Name(Name string) *UserQuery {
	p.db = p.db.Where(UserTable+".name = ?", Name)

	return p
}

func (p *UserQuery) NameLike(pattern string) *UserQuery {
	p.db = p.db.Where(UserTable+".name ilike ?", pattern)

	return p
}

func (p *UserQuery) NameIn(items []string) *UserQuery {
	p.db = p.db.Where(UserTable+".name IN (?)", items)

	return p
}

func (p *UserQuery) Email(Email string) *UserQuery {
	p.db = p.db.Where(UserTable+".email = ?", Email)

	return p
}

func (p *UserQuery) EmailLike(pattern string) *UserQuery {
	p.db = p.db.Where(UserTable+".email ilike ?", pattern)

	return p
}

func (p *UserQuery) EmailIn(items []string) *UserQuery {
	p.db = p.db.Where(UserTable+".email IN (?)", items)

	return p
}

func (p *UserQuery) RegisteredAt(RegisteredAt time.Time) *UserQuery {
	p.db = p.db.Where(UserTable+".registered_at = ?", RegisteredAt)

	return p
}

func (p *UserQuery) Build() *gorm.DB {
	return p.db
}
