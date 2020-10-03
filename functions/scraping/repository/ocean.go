package repository

import (
	"scraping/database"
	"scraping/logging"
	"scraping/model"
)

// OceanRepository ...
type OceanRepository interface {
	Store(ocean *model.Ocean) error
}

// OceanRepositoryImpl ...
type OceanRepositoryImpl struct {
	logging *logging.OceanLoggingImpl
	db      *database.DynamoDBServiceImpl
}

// NewOceanRepository ...
func NewOceanRepository(logging *logging.OceanLoggingImpl) *OceanRepositoryImpl {
	return &OceanRepositoryImpl{
		logging: logging,
		db:      database.NewDynamoDBService(),
	}
}

// Store ...
func (s *OceanRepositoryImpl) Store(ocean *model.Ocean) error {
	s.logging.Info("データの永久保存開始", ocean.LocationName)
	err := s.db.CreateIfNotExist(ocean)
	if err != nil {
		s.logging.Info("データの永久保存に失敗")
		return err
	}
	s.logging.Info("データの永久保存終了", ocean.LocationName)
	return nil
}
