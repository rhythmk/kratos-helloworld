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

func (r *studentRepo) Get(ctx context.Context, id int64) (*biz.Student, error) {
	var stu biz.Student
	//stu.Name = r.data.gormDB.
	r.data.db.Where("id=?", id).First(&stu)

	r.log.WithContext(ctx).Info("gormDB: GetStudent, id: ", id)
	return &stu, nil
}

func (r *studentRepo) Create(ctx context.Context, g *biz.Student) (*biz.Student, error) {
	result := r.data.db.Model(&biz.Student{}).Create(g)
	r.log.WithContext(ctx).Info("gormDB: Create, id: ", g.Id)
	return g, result.Error
}
