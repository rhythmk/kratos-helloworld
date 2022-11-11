package data

import (
	"context"

	"kratos-shop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewStudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *studentRepo) Get(ctx context.Context, g *biz.Student) (*biz.Student, error) {
	return g, nil
}

func (r *studentRepo) Create(ctx context.Context, g *biz.Student) (*biz.Student, error) {
	return g, nil
}
