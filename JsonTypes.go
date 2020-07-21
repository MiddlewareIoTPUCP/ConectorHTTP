package main

import "encoding/json"

type virtualModel struct {
	Type        string `json:"type" binding:"required"`
	ReadingType string `json:"readingType" binding:"required"`
	Units       string `json:"units" binding:"required"`
	DataType    string `json:"dataType" binding:"required"`
}

type newRegisterJSON struct {
	Operation      string         `json:"operation" binding:"required"`
	OwnerToken     string         `json:"ownerToken" binding:"required"`
	DeviceID       string         `json:"deviceID" binding:"required"`
	VirtualModel   []virtualModel `json:"virtualModel" binding:"required"`
	AdditionalInfo interface{}    `json:"additonalInfo,omitempty"`
}

type readingObj struct {
	Index   int         `json:"index" binding:"required"`
	Reading json.Number `json:"reading" binding:"required"`
}

type readingsJSON struct {
	ObjID      string       `json:"objID" binding:"required"`
	OwnerToken string       `json:"ownerToken" binding:"required"`
	Readings   []readingObj `json:"readings" binding:"required"`
}
