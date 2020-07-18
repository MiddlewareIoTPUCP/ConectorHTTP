package main

type virtualModel struct {
	Type        string `json:"type" binding:"required"`
	ReadingType string `json:"readingType" binding:"required"`
	Units       string `json:"units" binding:"required"`
	DataType    string `json:"dataType" binding:"required"`
}

type newRegister struct {
	Operation      string         `json:"operation" binding:"required"`
	OwnerToken     string         `json:"ownerToken" binding:"required"`
	DeviceID       string         `json:"deviceID" binding:"required"`
	VirtualModel   []virtualModel `json:"virtualModel" binding:"required"`
	AdditionalInfo interface{}    `json:"additonalInfo,omitempty"`
}
