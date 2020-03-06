package models

type Maps struct {
	Id       string        `json:"id"`
	Type     string        `json:"type"`
	Name     string        `json:"name"`
	Geofence [][][]float64 `json:"geofence"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type GetMaps struct {
	Id []string `json:"id"`
}

type Response_Get struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Maps
}

type Temp struct {
	Geofence string
}

type PutMaps struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Geofence [][][]float64 `json:"geofence"`
}

type DeleteMaps struct {
	Id []string `json:"id"`
}
type Check struct {
	Id string `json:"id"`
}