package service

import (
	"context"
	"errors"
	"fmt"

	pb "shortvid-backend/api/shortvid-service/v1"

	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type FileService struct {
	pb.UnimplementedFileServer
	logger log.Logger
}

func NewFileService(logger log.Logger) *FileService {
	return &FileService{logger: logger}
}

func (s *FileService) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileReply, error) {
	hr, ok := khttp.RequestFromServerContext(ctx)
	if !ok {
		return nil, errors.New("get http request failed")
	}

	if s.logger == nil {
		return nil, fmt.Errorf("file service logger is nil")
	}
	s.logger.Log(log.LevelInfo, "headers", hr.Header)
	return &pb.UploadFileReply{}, nil
}
func (s *FileService) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileReply, error) {
	return &pb.DeleteFileReply{}, nil
}
func (s *FileService) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileReply, error) {
	return &pb.GetFileReply{}, nil
}
func (s *FileService) ListFile(ctx context.Context, req *pb.ListFileRequest) (*pb.ListFileReply, error) {
	return &pb.ListFileReply{}, nil
}
