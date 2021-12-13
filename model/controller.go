package model

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Controller struct {
	gorm.Model
	CID   string `gorm:"uniqueIndex;column:cid" json:"cid"` // Controller ID
	CName string `gorm:"column:cname" json:"cname"`         // Device ID
	Key   string `gorm:"column:key" json:"key"`             // Service ID
}

func (s *dbHandler) AddController(r io.Reader) (*Controller, error) {
	decoder := json.NewDecoder(r)
	var controller = &Controller{}

	err := decoder.Decode(controller)
	if err != nil {
		return nil, err
	}

	controller.CID = uuid.NewString()
	controller.Key = controller.CID

	tx := s.db.Create(controller)
	if tx.Error != nil {
		return nil, tx.Error
	}

	tx.First(&controller, "cid=?", controller.CID)
	fmt.Println(controller)
	return controller, nil
}

func (s *dbHandler) GetControllers() ([]*Controller, error) {
	var list []*Controller

	result := s.db.Find(&list)

	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

func (s *dbHandler) IsExistController(cid string) bool {
	var controller = Controller{}

	result := s.db.First(&controller, "cid=?", cid)

	return result.Error == nil
}
