package provider

// RouteEntry는 라우팅 테이블 항목을 나타냅니다
type RouteEntry struct {
	Destination string // CIDR 또는 "default"
	Gateway     string // "" (게이트웨이 없는 경우)
	Interface   string // "en0", "eth0" 등
	Metric      int    // 메트릭 값 (optional, 0이면 표시하지 않음)
	Source      string // 소스 IP (optional)
}

// RouteProvider는 OS별 라우팅 테이블 조회 인터페이스입니다
type RouteProvider interface {
	GetRoutes() ([]RouteEntry, error)
}

