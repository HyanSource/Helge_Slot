package metrics

import (
	metrics "github.com/rcrowley/go-metrics"
)

func NewLocalMetrics() *LocalMetrics {
	return &LocalMetrics{}
}

type LocalMetrics struct {
	AllPlayCounter  metrics.Counter //總共的遊玩次數
	WinMoneyCounter metrics.Counter //贏的金錢
	RTPGauge        metrics.Gauge   //RTP
}
