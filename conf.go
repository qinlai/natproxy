package main

type Proxy struct {
	Port  uint16 `json:"port"`
	Proxy string `json:"proxy"`
}

type Config struct {
	GateWayAddr       string   `json:"gate_way_addr"`
	GateWayClientPort uint16   `json:"gate_way_client_port"`
	GateWayServerPort uint16   `json:"gate_way_server_port"`
	GateWayAuthKey    string   `json:"gate_way_auth_key"`
	Proxys            []*Proxy `json:"proxys"`
}
