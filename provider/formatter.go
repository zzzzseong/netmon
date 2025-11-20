package provider

import (
	"fmt"
	"strings"
)

// FormatLinuxStyle은 RouteEntry를 Linux ip route 스타일로 포맷팅합니다
func FormatLinuxStyle(entry RouteEntry) string {
	var parts []string

	// Destination
	parts = append(parts, entry.Destination)

	// via Gateway
	if entry.Gateway != "" {
		parts = append(parts, "via", entry.Gateway)
	}

	// dev Interface
	if entry.Interface != "" {
		parts = append(parts, "dev", entry.Interface)
	}

	// metric
	if entry.Metric > 0 {
		parts = append(parts, "metric", fmt.Sprintf("%d", entry.Metric))
	}

	// src
	if entry.Source != "" {
		parts = append(parts, "src", entry.Source)
	}

	return strings.Join(parts, " ")
}

// FormatLinuxStyleRoutes는 RouteEntry 슬라이스를 Linux ip route 스타일로 포맷팅합니다
func FormatLinuxStyleRoutes(entries []RouteEntry) string {
	var lines []string
	for _, entry := range entries {
		lines = append(lines, FormatLinuxStyle(entry))
	}
	return strings.Join(lines, "\n")
}

