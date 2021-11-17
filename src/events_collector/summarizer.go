package events_collector

type Summarizer interface {
	Summarize(metricsPromise CollectedMetricsPromise) (string, bool)
}
