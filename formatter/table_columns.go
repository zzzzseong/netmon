package formatter

// Table column width constants for port table
const (
	PortTableProtocolWidth     = 10
	PortTableLocalAddressWidth = 19
	PortTableStatusWidth       = 10
	PortTablePIDWidth          = 8
	PortTableProcessNameWidth  = 25
	PortTableUsernameWidth     = 15
	PortTableCPUWidth          = 8
	PortTableMemWidth          = 8
)

// PortTableColumns defines the column configuration for the port table
var PortTableColumns = []TableColumn{
	{Width: PortTableProtocolWidth, Title: "PROTOCOL"},
	{Width: PortTableLocalAddressWidth, Title: "LOCAL ADDRESS"},
	{Width: PortTableStatusWidth, Title: "STATUS"},
	{Width: PortTablePIDWidth, Title: "PID"},
	{Width: PortTableProcessNameWidth, Title: "PROCESS NAME"},
	{Width: PortTableUsernameWidth, Title: "USERNAME"},
	{Width: PortTableCPUWidth, Title: "CPU %"},
	{Width: PortTableMemWidth, Title: "MEM %"},
}

// Table column width constants for route table
const (
	RouteTableDestinationWidth = 27
	RouteTableGatewayWidth     = 27
	RouteTableInterfaceWidth   = 27
	RouteTableMetricWidth      = 27
	RouteTableSourceWidth      = 26
)

// RouteTableColumns defines the column configuration for the route table
var RouteTableColumns = []TableColumn{
	{Width: RouteTableDestinationWidth, Title: "DESTINATION"},
	{Width: RouteTableGatewayWidth, Title: "GATEWAY"},
	{Width: RouteTableInterfaceWidth, Title: "INTERFACE"},
	{Width: RouteTableMetricWidth, Title: "METRIC"},
	{Width: RouteTableSourceWidth, Title: "SOURCE"},
}

// Table column width constants for interface table
const (
	InterfaceTableNameWidth      = 14
	InterfaceTableIPAddressWidth = 42
	InterfaceTableMACAddressWidth = 22
	InterfaceTableStatusWidth    = 13
	InterfaceTableMTUWidth       = 13
)

// InterfaceTableColumns defines the column configuration for the interface table
var InterfaceTableColumns = []TableColumn{
	{Width: InterfaceTableNameWidth, Title: "INTERFACE"},
	{Width: InterfaceTableIPAddressWidth, Title: "IP ADDRESS"},
	{Width: InterfaceTableMACAddressWidth, Title: "MAC ADDRESS"},
	{Width: InterfaceTableStatusWidth, Title: "STATUS"},
	{Width: InterfaceTableMTUWidth, Title: "MTU"},
}

// Table column width constants for traceroute table
const (
	TracerouteTableHopWidth  = 10
	TracerouteTableHostWidth = 42
	TracerouteTableRTT1Width = 14
	TracerouteTableRTT2Width = 14
	TracerouteTableRTT3Width = 14
)

// TracerouteTableColumns defines the column configuration for the traceroute table
var TracerouteTableColumns = []TableColumn{
	{Width: TracerouteTableHopWidth, Title: "HOP"},
	{Width: TracerouteTableHostWidth, Title: "HOST"},
	{Width: TracerouteTableRTT1Width, Title: "RTT 1"},
	{Width: TracerouteTableRTT2Width, Title: "RTT 2"},
	{Width: TracerouteTableRTT3Width, Title: "RTT 3"},
}

// Table column width constants for connection table
const (
	ConnectionTableProtocolWidth      = 10
	ConnectionTableLocalAddressWidth  = 25
	ConnectionTableRemoteAddressWidth = 25
	ConnectionTablePIDWidth           = 10
	ConnectionTableProcessNameWidth   = 25
)

// ConnectionTableColumns defines the column configuration for the connection table
var ConnectionTableColumns = []TableColumn{
	{Width: ConnectionTableProtocolWidth, Title: "PROTOCOL"},
	{Width: ConnectionTableLocalAddressWidth, Title: "LOCAL ADDRESS"},
	{Width: ConnectionTableRemoteAddressWidth, Title: "REMOTE ADDRESS"},
	{Width: ConnectionTablePIDWidth, Title: "PID"},
	{Width: ConnectionTableProcessNameWidth, Title: "PROCESS"},
}

// Table column width constants for DNS table
const (
	DNSTableRecordTypeWidth = 12
	DNSTableValueWidth      = 60
)

// DNSTableColumns defines the column configuration for the DNS table
var DNSTableColumns = []TableColumn{
	{Width: DNSTableRecordTypeWidth, Title: "TYPE"},
	{Width: DNSTableValueWidth, Title: "VALUE"},
}


