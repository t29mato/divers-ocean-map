package service

import (
	"scraping/logging"
	"scraping/model"
)

// RepositoryService ...
type RepositoryService interface {
	Store(ocean *model.Ocean) error
}

// RepositoryServiceImpl ...
type RepositoryServiceImpl struct {
	logging *logging.OceanLoggingImpl
	db      *DynamoDBServiceImpl
}

// NewRepository ...
func NewRepository(logging *logging.OceanLoggingImpl) *RepositoryServiceImpl {
	return &RepositoryServiceImpl{
		logging: logging,
		db:      NewDynamoDBService(),
	}
}

// Store ...
func (s *RepositoryServiceImpl) Store(ocean *model.Ocean) error {
	s.logging.Info("データの永久保存開始", ocean.LocationName)
	err := s.db.CreateIfNotExist(ocean)
	if err != nil {
		s.logging.Info("データの永久保存に失敗")
		return err
	}
	s.logging.Info("データの永久保存終了", ocean.LocationName)
	return nil
}
