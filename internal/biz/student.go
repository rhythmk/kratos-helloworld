package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Student struct {
	Id   uint64 `gorm:"column:id"                json:"id"                 compare:"Id"`
	Name string `gorm:"column:name"                json:"name"                 compare:"Name"`
}

func (t *Student) TableName() string {
	return "student"
}

type StudentRepo interface {
	Get(context.Context, int64) (*Student, error)
	Create(context.Context, *Student) (*Student, error)
}

type StudentUsercase struct {
	repo StudentRepo
	log  *log.Helper
}

func NewStudentUsercase(repo StudentRepo, logger log.Logger) *StudentUsercase {
	return &StudentUsercase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *StudentUsercase) CreateStudent(ctx context.Context, stu *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("CreateStudent: %v", stu.Id)
	return uc.repo.Create(ctx, stu)
}

func (uc *StudentUsercase) GetStudent(ctx context.Context, stu *Student) (*Student, error) {
	uc.log.WithContext(ctx).Infof("GetStudent: %v", stu.Id)
	return uc.repo.Get(ctx, int64(stu.Id))
}
