package service

import (
	"context"
	"fmt"
	pb "kratos-shop/api/helloworld/v1"
	"kratos-shop/internal/biz"
)

type StudentService struct {
	pb.UnimplementedStudentServer
	uc *biz.StudentUsercase
}

func NewStudentService(uc *biz.StudentUsercase) *StudentService {
	return &StudentService{uc: uc}
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {

	r, err := s.uc.CreateStudent(ctx, &biz.Student{
		Id:   req.Id,
		Name: req.Name,
	})
	return &pb.CreateStudentReply{
		Code: "创建成功",
		Msg:  fmt.Sprintf("id:%v,name:%v,err:%v", r.Id, r.Name, err.Error()),
	}, nil
}
func (s *StudentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	r, err := s.GetStudent(ctx, &pb.GetStudentRequest{
		Id: req.Id,
	})
	if err != nil {
		panic(err)
	}
	return &pb.GetStudentReply{
		Id:   r.Id,
		Name: r.Name,
	}, nil
}
