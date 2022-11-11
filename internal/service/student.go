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
	return &pb.CreateStudentReply{
		Code: "创建成功",
		Msg:  fmt.Sprintf("id:%v,name:%v", req.Id, req.Name),
	}, nil
}
func (s *StudentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	return &pb.GetStudentReply{
		Id:   req.Id,
		Name: "用户名",
	}, nil
}
