package models

type Valve struct {
	ID           	uint    `json:"id" gorm:"primarykey"`
	Manufacturer 	string  `json:"manufacturer"`
	PartNumber   	string  `json:"part_number" gorm:"unique;not null"`
	ComponentSeries string  `json:"component_series" gorm:"not null"`
	ValveType    	string  `json:"valve_type" gorm:"not null"`    
	Size_NG			int     `json:"size_ng"`    
	Weight      	float64 `json:"weight"`
	RatedFlow     	float64 `json:"rated_flow"`                   
	MaxFlow   		float64 `json:"max_flow"`
	CommandValue    float64 `json:"command_value"`
	CommandType     string  `json:"command_type"`
	MaxPressure  	float64 `json:"max_pressure"`
	ImagePath       string  `json:"image_path"`   
	DatasheetPath   string  `json:"datasheet_path"`                       
}