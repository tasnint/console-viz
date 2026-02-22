// Package collector fetches and parses metrics from exporters (e.g. Windows Exporter).
package collector

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CoreFrequency is one logical core's current frequency in MHz (for windows_cpu_core_frequency_mhz).
type CoreFrequency struct {
	// Core is the core id from the metric label, e.g. "0,0", "0,1"
	Core string
	// Mhz is the frequency in MHz at scrape time
	Mhz float64
}

// CPUFrequencySnapshot is a point-in-time snapshot of all core frequencies (one sample per 60s fetch).
type CPUFrequencySnapshot struct {
	// Time is when we fetched the metrics (x-axis for time series)
	Time time.Time
	// Cores lists each core's frequency so we can plot or aggregate
	Cores []CoreFrequency
}

// FetchCPUFrequency fetches metrics from the given URL, parses Prometheus text format,
// and returns the current windows_cpu_core_frequency_mhz values. Call this every 60s from main's ticker.
func FetchCPUFrequency(metricsURL string) (*CPUFrequencySnapshot, error) {
	// record when we fetched so x-axis (time) always increases
	now := time.Now()
	// HTTP GET the metrics endpoint (e.g. http://localhost:9182/metrics)
	resp, err := http.Get(metricsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch metrics: %w", err)
	}
	// ensure we close the body when we're done so we don't leak connections
	defer resp.Body.Close()
	// only accept 200 OK; otherwise body might be an error page
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch metrics: status %s", resp.Status)
	}
	// read the whole body so we can split by lines (Prometheus format is text)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read metrics body: %w", err)
	}
	text := string(body)
	// parse the text and extract only windows_cpu_core_frequency_mhz samples
	cores, err := parseCPUFrequencyMetric(text)
	if err != nil {
		return nil, fmt.Errorf("parse metrics: %w", err)
	}
	// return one snapshot (time + values) for this 60s interval
	return &CPUFrequencySnapshot{Time: now, Cores: cores}, nil
}

// coresToTrack lists which cores we include; each will be a separate line on the graph later.
// For now only core "0,0"; add e.g. "0,1", "0,2" to get more lines on the same graph.
var coresToTrack = []string{"0,0", "0,1", "0,2"}

// parseCPUFrequencyMetric finds lines like: windows_cpu_core_frequency_mhz{core="0,0"} 1506
// and returns a slice of CoreFrequency. Skips # comments and non-matching metrics.
// Only includes cores in coresToTrack so we track one line first, then add more as needed.
func parseCPUFrequencyMetric(text string) ([]CoreFrequency, error) {
	// we'll append each parsed core here (only if it's in coresToTrack)
	var out []CoreFrequency
	// Prometheus format is one metric per line; split so we can iterate
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		// skip empty lines and comment lines (# HELP, # TYPE, or #)
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// we only care about this metric name
		if !strings.HasPrefix(line, "windows_cpu_core_frequency_mhz") {
			continue
		}
		// line shape: windows_cpu_core_frequency_mhz{core="0,0"} 1506
		// find the value: last field after whitespace
		idx := strings.LastIndex(line, " ")
		if idx == -1 {
			continue
		}
		valueStr := strings.TrimSpace(line[idx:])
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}
		// extract label core="x,y" from the { ... } part
		start := strings.Index(line, "{")
		end := strings.Index(line, "}")
		core := ""
		if start != -1 && end != -1 && end > start {
			labels := line[start+1 : end]
			// simple parse: look for core="..."
			if strings.HasPrefix(labels, "core=\"") {
				closeQuote := strings.Index(labels[6:], "\"")
				if closeQuote != -1 {
					core = labels[6 : 6+closeQuote]
				}
			}
		}
		// include only cores we want to track (e.g. "0,0" for now; add more later for separate lines)
		keep := false
		for _, c := range coresToTrack {
			if core == c {
				keep = true
				break
			}
		}
		if !keep {
			continue
		}
		// append this core's data for this snapshot
		out = append(out, CoreFrequency{Core: core, Mhz: value})
	}
	return out, nil
}
