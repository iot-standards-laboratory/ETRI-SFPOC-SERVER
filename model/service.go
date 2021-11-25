package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	SName     string `gorm:"column:sname" json:"sname"`
	SID       string `gorm:"column:sid" json:"sid"`
	NumOfDevs int    `gorm:"column:ndevs" json:"ndevs"`
	Addr      string `gorm:"column:addr"`
}

func (s *dbHandler) GetServices() ([]*Service, error) {
	var services []*Service

	result := s.db.Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	// var devices []*Device
	for _, service := range services {
		tx := s.db.Where("sname=?", service.SName).Find(&[]*Device{})
		if tx.Error == gorm.ErrRecordNotFound {
			service.NumOfDevs = 0
		} else if tx.Error != nil {
			return nil, tx.Error
		} else {
			service.NumOfDevs = int(tx.RowsAffected)
		}
	}
	return services, nil
}

func (s *dbHandler) AddService(name string) error {
	result := s.db.First(&Service{}, "sname=?", name)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			tx := s.db.Create(&Service{SName: name})
			if tx.Error != nil {
				return tx.Error
			}
		}

		return result.Error
	}

	return nil
}

func (s *dbHandler) IsExistService(name string) bool {
	var service Service
	tx := s.db.Select("sid").First(&service, "sname=?", name)

	return tx.Error == nil
}

func (s *dbHandler) UpdateService(name, addr string) (*Service, error) {
	tx := s.db.Model(&Service{}).Where("sname = ?", name).Updates(Service{SName: name, SID: uuid.NewString(), Addr: addr})
	if tx.Error != nil {
		return nil, tx.Error
	}

	var service Service
	tx.First(&service, "sname=?", name)
	return &service, nil
}

func (s *dbHandler) GetSID(name string) (string, error) {
	var service Service
	tx := s.db.Select("sid").First(&service, "sname=?", name)
	if tx.Error != nil {
		return "", tx.Error
	}

	return service.SID, nil
}

func (s *dbHandler) GetAddr(sid string) (string, error) {
	var service Service

	result := s.db.First(&service, "sid=?", sid)

	if result.Error != nil {
		return "", result.Error
	}

	return service.Addr, nil
}
