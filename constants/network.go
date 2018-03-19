package constants

import "Go-Ground-Station/gstypes"

var (
	Hosts = map[int]gstypes.Node{
		9110: {"192.168.0.255", 9110, "Power Node A"},
		9111: {"192.168.0.255", 9111, "Power Node B"},
		9531: {"192.168.0.255", 9531, "Flight Control"},
		9120: {"192.168.0.255", 9120, "Landing Gear"},
		9130: {"192.168.0.255", 9130, "Gimbal Control"},
		9170: {"192.168.0.255", 9170, "Xilinx Sim"},
	}

	GrpcPort = 9800
)
