package collector

import (
	"fmt"

	"github.com/genma/nzbget-metrics/internal/nzbget"
	"github.com/prometheus/client_golang/prometheus"
)

type nzbgetClient interface {
	Status() (*nzbget.StatusResult, error)
	Version() (string, error)
	QueuedCount() (int, error)
}

type Collector struct {
	client nzbgetClient

	up                        *prometheus.Desc
	info                      *prometheus.Desc
	downloadRate              *prometheus.Desc
	averageDownloadRate       *prometheus.Desc
	downloadLimit             *prometheus.Desc
	remainingBytes            *prometheus.Desc
	downloadedBytes           *prometheus.Desc
	articleCacheBytes         *prometheus.Desc
	freeDiskBytes             *prometheus.Desc
	totalDiskBytes            *prometheus.Desc
	uptimeSeconds             *prometheus.Desc
	downloadTimeSeconds       *prometheus.Desc
	threadCount               *prometheus.Desc
	postJobCount              *prometheus.Desc
	urlCount                  *prometheus.Desc
	queueScriptCount          *prometheus.Desc
	queuedCount               *prometheus.Desc
	paused                    *prometheus.Desc
	serverStandby             *prometheus.Desc
	quotaReached              *prometheus.Desc
	daySizeBytes              *prometheus.Desc
	monthSizeBytes            *prometheus.Desc
	newsServerActive          *prometheus.Desc
}

func New(client nzbgetClient) *Collector {
	return &Collector{
		client: client,
		up:                  prometheus.NewDesc("nzbget_up", "1 if the nzbget API is reachable.", nil, nil),
		info:                prometheus.NewDesc("nzbget_info", "nzbget version info.", []string{"version"}, nil),
		downloadRate:        prometheus.NewDesc("nzbget_download_rate_bytes_per_second", "Current download rate in bytes/s.", nil, nil),
		averageDownloadRate: prometheus.NewDesc("nzbget_average_download_rate_bytes_per_second", "Average download rate in bytes/s.", nil, nil),
		downloadLimit:       prometheus.NewDesc("nzbget_download_limit_bytes_per_second", "Download speed limit in bytes/s.", nil, nil),
		remainingBytes:      prometheus.NewDesc("nzbget_remaining_size_bytes", "Remaining download size in bytes.", nil, nil),
		downloadedBytes:     prometheus.NewDesc("nzbget_downloaded_size_bytes", "Total downloaded size in bytes.", nil, nil),
		articleCacheBytes:   prometheus.NewDesc("nzbget_article_cache_bytes", "Article cache usage in bytes.", nil, nil),
		freeDiskBytes:       prometheus.NewDesc("nzbget_free_disk_space_bytes", "Free disk space in bytes.", nil, nil),
		totalDiskBytes:      prometheus.NewDesc("nzbget_total_disk_space_bytes", "Total disk space in bytes.", nil, nil),
		uptimeSeconds:       prometheus.NewDesc("nzbget_uptime_seconds", "nzbget uptime in seconds.", nil, nil),
		downloadTimeSeconds: prometheus.NewDesc("nzbget_download_time_seconds", "Total download time in seconds.", nil, nil),
		threadCount:         prometheus.NewDesc("nzbget_threads", "Number of active threads.", nil, nil),
		postJobCount:        prometheus.NewDesc("nzbget_post_jobs", "Number of post-processing jobs.", nil, nil),
		urlCount:            prometheus.NewDesc("nzbget_urls", "Number of URL fetches in queue.", nil, nil),
		queueScriptCount:    prometheus.NewDesc("nzbget_queue_scripts", "Number of queued scripts.", nil, nil),
		queuedCount:         prometheus.NewDesc("nzbget_queued", "Number of NZBs in queue.", nil, nil),
		paused:              prometheus.NewDesc("nzbget_paused", "1 if the given component is paused.", []string{"kind"}, nil),
		serverStandby:       prometheus.NewDesc("nzbget_server_standby", "1 if nzbget is in standby (no active downloads).", nil, nil),
		quotaReached:        prometheus.NewDesc("nzbget_quota_reached", "1 if the download quota has been reached.", nil, nil),
		daySizeBytes:        prometheus.NewDesc("nzbget_day_size_bytes", "Data downloaded today in bytes.", nil, nil),
		monthSizeBytes:      prometheus.NewDesc("nzbget_month_size_bytes", "Data downloaded this month in bytes.", nil, nil),
		newsServerActive:    prometheus.NewDesc("nzbget_news_server_active", "1 if the news server is active.", []string{"id"}, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.downloadRate
	ch <- c.averageDownloadRate
	ch <- c.downloadLimit
	ch <- c.remainingBytes
	ch <- c.downloadedBytes
	ch <- c.articleCacheBytes
	ch <- c.freeDiskBytes
	ch <- c.totalDiskBytes
	ch <- c.uptimeSeconds
	ch <- c.downloadTimeSeconds
	ch <- c.threadCount
	ch <- c.postJobCount
	ch <- c.urlCount
	ch <- c.queueScriptCount
	ch <- c.queuedCount
	ch <- c.paused
	ch <- c.serverStandby
	ch <- c.quotaReached
	ch <- c.daySizeBytes
	ch <- c.monthSizeBytes
	ch <- c.newsServerActive
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	status, err := c.client.Status()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	queued, err := c.client.QueuedCount()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

	version, _ := c.client.Version()
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, version)

	ch <- prometheus.MustNewConstMetric(c.downloadRate, prometheus.GaugeValue, status.DownloadRate)
	ch <- prometheus.MustNewConstMetric(c.averageDownloadRate, prometheus.GaugeValue, status.AverageDownloadRate)
	ch <- prometheus.MustNewConstMetric(c.downloadLimit, prometheus.GaugeValue, status.DownloadLimit)

	ch <- prometheus.MustNewConstMetric(c.remainingBytes, prometheus.GaugeValue, status.RemainingBytes())
	ch <- prometheus.MustNewConstMetric(c.downloadedBytes, prometheus.GaugeValue, status.DownloadedBytes())
	ch <- prometheus.MustNewConstMetric(c.articleCacheBytes, prometheus.GaugeValue, status.ArticleCacheBytes())
	ch <- prometheus.MustNewConstMetric(c.freeDiskBytes, prometheus.GaugeValue, status.FreeDiskBytes())
	ch <- prometheus.MustNewConstMetric(c.totalDiskBytes, prometheus.GaugeValue, status.TotalDiskBytes())

	ch <- prometheus.MustNewConstMetric(c.uptimeSeconds, prometheus.GaugeValue, float64(status.UpTimeSec))
	ch <- prometheus.MustNewConstMetric(c.downloadTimeSeconds, prometheus.GaugeValue, float64(status.DownloadTimeSec))
	ch <- prometheus.MustNewConstMetric(c.threadCount, prometheus.GaugeValue, float64(status.ThreadCount))
	ch <- prometheus.MustNewConstMetric(c.postJobCount, prometheus.GaugeValue, float64(status.PostJobCount))
	ch <- prometheus.MustNewConstMetric(c.urlCount, prometheus.GaugeValue, float64(status.UrlCount))
	ch <- prometheus.MustNewConstMetric(c.queueScriptCount, prometheus.GaugeValue, float64(status.QueueScriptCount))
	ch <- prometheus.MustNewConstMetric(c.queuedCount, prometheus.GaugeValue, float64(queued))

	boolGauge := func(b bool) float64 {
		if b {
			return 1
		}
		return 0
	}

	ch <- prometheus.MustNewConstMetric(c.paused, prometheus.GaugeValue, boolGauge(status.DownloadPaused), "download")
	ch <- prometheus.MustNewConstMetric(c.paused, prometheus.GaugeValue, boolGauge(status.ServerPaused), "server")
	ch <- prometheus.MustNewConstMetric(c.paused, prometheus.GaugeValue, boolGauge(status.PostPaused), "post")
	ch <- prometheus.MustNewConstMetric(c.paused, prometheus.GaugeValue, boolGauge(status.ScanPaused), "scan")
	ch <- prometheus.MustNewConstMetric(c.serverStandby, prometheus.GaugeValue, boolGauge(status.ServerStandBy))
	ch <- prometheus.MustNewConstMetric(c.quotaReached, prometheus.GaugeValue, boolGauge(status.QuotaReached))

	ch <- prometheus.MustNewConstMetric(c.daySizeBytes, prometheus.GaugeValue, status.DayBytes())
	ch <- prometheus.MustNewConstMetric(c.monthSizeBytes, prometheus.GaugeValue, status.MonthBytes())

	for _, ns := range status.NewsServers {
		id := fmt.Sprintf("%d", ns.ID)
		ch <- prometheus.MustNewConstMetric(c.newsServerActive, prometheus.GaugeValue, boolGauge(ns.Active), id)
	}
}
