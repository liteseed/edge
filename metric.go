package bungo

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/big"
)

const (
	MetricNameSpace = "bungo"
)

var (
	bundlerBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: MetricNameSpace,
			Name:      "bundler_balance",
			Help:      "upload data to arweave",
		},
		[]string{"bundler", "token"},
	)
)

func init() {
	prometheus.MustRegister(
		bundlerBalance,
	)
}

func metricBundlerBalance(bal *big.Float, addr string) {
	amount, _ := bal.Float64()
	bundlerBalance.WithLabelValues(addr, "AR").Set(amount)
}
