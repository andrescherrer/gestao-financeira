package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestDuration tracks HTTP request duration in seconds
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestTotal tracks total number of HTTP requests
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestInFlight tracks number of in-flight HTTP requests
	HTTPRequestInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of in-flight HTTP requests",
		},
	)

	// BusinessMetrics tracks business-specific metrics
	BusinessMetrics = struct {
		TransactionsCreated *prometheus.CounterVec
		TransactionsUpdated *prometheus.CounterVec
		TransactionsDeleted *prometheus.CounterVec
		AccountsCreated     *prometheus.CounterVec
		AccountsUpdated     *prometheus.CounterVec
		CategoriesCreated   *prometheus.CounterVec
		BudgetsCreated      *prometheus.CounterVec
		InvestmentsCreated  *prometheus.CounterVec
		UsersRegistered     prometheus.Counter
		LoginAttempts       *prometheus.CounterVec
	}{
		TransactionsCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_transactions_created_total",
				Help: "Total number of transactions created",
			},
			[]string{"type"},
		),
		TransactionsUpdated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_transactions_updated_total",
				Help: "Total number of transactions updated",
			},
			[]string{"type"},
		),
		TransactionsDeleted: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_transactions_deleted_total",
				Help: "Total number of transactions deleted",
			},
			[]string{"type"},
		),
		AccountsCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_accounts_created_total",
				Help: "Total number of accounts created",
			},
			[]string{"type"},
		),
		AccountsUpdated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_accounts_updated_total",
				Help: "Total number of accounts updated",
			},
			[]string{"type"},
		),
		CategoriesCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_categories_created_total",
				Help: "Total number of categories created",
			},
			[]string{"type"},
		),
		BudgetsCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_budgets_created_total",
				Help: "Total number of budgets created",
			},
			[]string{"period_type"},
		),
		InvestmentsCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_investments_created_total",
				Help: "Total number of investments created",
			},
			[]string{"type"},
		),
		UsersRegistered: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "business_users_registered_total",
				Help: "Total number of users registered",
			},
		),
		LoginAttempts: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_login_attempts_total",
				Help: "Total number of login attempts",
			},
			[]string{"status"},
		),
	}
)

// Init initializes Prometheus metrics
func Init() {
	// Metrics are automatically registered via promauto
	// This function can be used for any additional initialization if needed
}
