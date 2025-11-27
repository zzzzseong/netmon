package main

import (
	"netmon/cmd"
)

// Version 정보
const (
	Version   = "1.1.5"
	BuildDate = "2025-11-26"
)

func main() {
	// cmd 패키지에 버전 정보 설정
	cmd.Version = Version
	cmd.BuildDate = BuildDate

	// Cobra 명령어 실행
	cmd.Execute()
}
