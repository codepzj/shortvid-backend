package data

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/data/model"
)

type videoRepo struct {
	data *Data
}

func NewVideoRepo(data *Data) *videoRepo {
	return &videoRepo{data: data}
}

func (r *videoRepo) CreateVideo(ctx context.Context, video *model.Video) error {
	return r.data.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepo) GetVideoByID(ctx context.Context, id int) (*model.Video, error) {
	var video model.Video
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}