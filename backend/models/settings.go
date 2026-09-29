package models

type Settings struct {
	ID                  uint   `json:"id" gorm:"primaryKey"`
	OpcUaAddress        string `json:"opc_ua_address"`
	OpcNodePressure     string `json:"opc_node_pressure"`
	OpcNodeCommand      string `json:"opc_node_command"`
	OpcNodeFeedback     string `json:"opc_node_feedback"`
	OpcNodeFlow         string `json:"opc_node_flow"`
	OpcNodeOutput       string `json:"opc_node_output"`
	OpcUaUsername       string `json:"opc_ua_username"`
	OpcUaPassword       string `json:"-"`
	OpcUaSecurityPolicy string `json:"opc_ua_security_policy"`
	OpcUaSecurityMode   string `json:"opc_ua_security_mode"`
}