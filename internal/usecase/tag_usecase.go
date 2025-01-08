package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/model/converter"
	"github.com/abdisetiakawan/go-clean-arch/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TagUseCase struct {
	DB *gorm.DB
	Log *logrus.Logger
	Validate *validator.Validate
	TagRepository *repository.TagRepository
	Cache *helper.CacheHelper
}

func NewTagUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, tagRepository *repository.TagRepository, cache *helper.CacheHelper) *TagUseCase {
	return &TagUseCase{
		DB: db,
		Log: log,
		Validate: validate,
		TagRepository: tagRepository,
		Cache: cache,
	}
}

func (c *TagUseCase) Create(ctx context.Context, request *model.CreateTagRequest) (*model.TagResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request body")
		return nil, model.ErrBadRequest
	}
	tag := &entity.Tag{
		Email: request.Email,
		Name: request.Name,
	}
	if err := c.TagRepository.Create(tx, tag); err != nil {
		c.Log.WithError(err).Error("error create tag")
		return nil, model.ErrInternalServer
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error create tag")
		return nil, model.ErrInternalServer
	}

	return converter.TagToResponse(tag), nil
}

func (c *TagUseCase) Search(ctx context.Context, request *model.SearchTagRequest) ([]model.TagResponse, int64, error) {
	cacheKey := "tags:search:" + request.Email
	var cachedData struct {
		Responses []model.TagResponse
		Total int64
	}
	if err := c.Cache.GetAndUnmarshal(ctx, cacheKey, &cachedData); err == nil {
		return cachedData.Responses, cachedData.Total, nil
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request body")
		return nil, 0, model.ErrBadRequest
	}
	tags, total, err := c.TagRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error search tag")
		return nil, 0, model.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error search tag")
		return nil, 0, model.ErrInternalServer
	}
	responses := make([]model.TagResponse, len(tags))
	for i, tag := range tags {
		tag.Email = ""
		responses[i] = *converter.TagToResponse(&tag)
	} 

	cachedData.Responses = responses
	cachedData.Total = total
	cachedDataJSON, _ := json.Marshal(cachedData)
	c.Cache.Set(ctx, cacheKey, cachedDataJSON, 1*time.Minute)
	
	return responses, total, nil
}

func (c *TagUseCase) Get(ctx context.Context, request *model.GetTagRequest) (*model.TagResponse, error) {
	var tagResponse model.TagResponse
	cacheKey := "tags:" + request.ID
	if err := c.Cache.GetAndUnmarshal(ctx, cacheKey, &tagResponse); err == nil {
		return &tagResponse, nil
	}
	
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request query")
		return nil, model.ErrBadRequest
	}
	tag := new(entity.Tag)
	if err := c.TagRepository.FindByEmailAndId(tx, tag, request.ID, request.Email); err != nil {
		c.Log.WithError(err).Error("error search tag")
		return nil, model.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error search tag")
		return nil, model.ErrInternalServer
	}

	tagResponse = *converter.TagToResponse(tag)
	tagResponseJSON, _ := json.Marshal(tagResponse)
	c.Cache.Set(ctx, cacheKey, tagResponseJSON, 1*time.Minute)
	
	return &tagResponse, nil
}

func (c *TagUseCase) Update(ctx context.Context, request *model.UpdateTagRequest) (*model.TagResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	tag := new(entity.Tag)
	if err := c.TagRepository.FindByEmailAndId(tx, tag, request.ID, request.Email); err != nil {
		c.Log.WithError(err).Error("error search tag")
		return nil, model.ErrNotFound
	}
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request query")
		return nil, model.ErrBadRequest
	}
	if request.Name != "" {
		tag.Name = request.Name
	}
	
	if err := c.TagRepository.Update(tx, tag); err != nil {
		c.Log.WithError(err).Error("error update tag")
		return nil, model.ErrInternalServer
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error update tag")
		return nil, model.ErrInternalServer
	}

	tagResponse := converter.TagToResponse(tag)
	tagResponseJSON, _ := json.Marshal(tagResponse)
	c.Cache.Set(ctx, "tags:"+request.ID, tagResponseJSON, 1*time.Minute)
	
	return tagResponse, nil
}

func (c *TagUseCase) Delete(ctx context.Context, request *model.GetTagRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request query")
		return model.ErrBadRequest
	}

	tag := new(entity.Tag)
	if err := c.TagRepository.FindByEmailAndId(tx, tag, request.ID, request.Email); err != nil {
		c.Log.WithError(err).Error("error search tag")
		return model.ErrNotFound
	}

	if err := c.TagRepository.Delete(tx, tag); err != nil {
		c.Log.WithError(err).Error("error delete tag")
		return model.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error delete tag")
		return model.ErrInternalServer
	}

	c.Cache.Delete(ctx, "tags:"+request.ID)

	return nil
}