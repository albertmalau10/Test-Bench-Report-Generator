package models

import "time"

type Settings struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	ModbusServerAddress  string     `json:"modbus_server_address"`
	CtrlxHost            string     `json:"ctrlx_host"`
	CtrlxUsername        string     `json:"ctrlx_username"`
	CtrlxPassword        string     `json:"-"` 
	CtrlxToken           string     `json:"ctrlx_token,omitempty"`
	CtrlxTokenExpiresAt  *time.Time `json:"ctrlx_token_expires_at,omitempty"`
}