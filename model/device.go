package model

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	DID   string `gorm:"uniqueIndex;column:did" json:"did"` // Device ID
	DName string `gorm:"column:dname" json:"dname"`         // Device Name
	Type  string `gorm:"column:type" json:"type"`           // Device Type
	CID   string `gorm:"column:cid" json:"cid"`             // Controller ID
	SID   string `gorm:"column:sid" json:"sid"`             // Service ID
	SName string `gorm:"column:sname" json:"sname"`         // Service Name
	// Opts []string
}

func (s *dbHandler) GetDevices() ([]*Device, int, error) {
	var devices []*Device

	result := s.db.Find(&devices)

	if result.Error != nil {
		return nil, -1, result.Error
	}
	return devices, int(result.RowsAffected), nil
}

// func (s *dbHandler) GetDevice() *Device {
// 	device := &Device{}
// }

func (s *dbHandler) AddDevice(device *Device) error {

	tx := s.db.Create(device)
	if tx.Error != nil {
		return tx.Error
	}

	tx.First(device, "did=?", device.DID)
	return nil

}

func (s *dbHandler) QueryDevice(dname string) (*Device, error) {
	var device Device
	tx := s.db.First(&device, "dname=?", dname)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &device, nil
}

func (s *dbHandler) DeleteDevice(device *Device) error {
	tx := s.db.Delete(device)
	if tx.Error != nil {
		return tx.Error
	}

	// tx.First(device, "did=?", device.DID)
	return nil
}

func (s *dbHandler) IsExistDevice(dname string) bool {
	var device = Device{}

	result := s.db.First(&device, "dname=?", dname)

	return result.Error != nil
}
