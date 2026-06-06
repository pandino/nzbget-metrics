package nzbget

type StatusResult struct {
	RemainingSizeMB       float64 `json:"RemainingSizeMB"`
	RemainingSizeLo       int64   `json:"RemainingSizeLo"`
	RemainingSizeHi       int64   `json:"RemainingSizeHi"`
	DownloadedSizeMB      float64 `json:"DownloadedSizeMB"`
	DownloadedSizeLo      int64   `json:"DownloadedSizeLo"`
	DownloadedSizeHi      int64   `json:"DownloadedSizeHi"`
	ArticleCacheMB        float64 `json:"ArticleCacheMB"`
	ArticleCacheLo        int64   `json:"ArticleCacheLo"`
	ArticleCacheHi        int64   `json:"ArticleCacheHi"`
	FreeDiskSpaceMB       float64 `json:"FreeDiskSpaceMB"`
	FreeDiskSpaceLo       int64   `json:"FreeDiskSpaceLo"`
	FreeDiskSpaceHi       int64   `json:"FreeDiskSpaceHi"`
	TotalDiskSpaceMB      float64 `json:"TotalDiskSpaceMB"`
	TotalDiskSpaceLo      int64   `json:"TotalDiskSpaceLo"`
	TotalDiskSpaceHi      int64   `json:"TotalDiskSpaceHi"`
	DaySizeMB             float64 `json:"DaySizeMB"`
	DaySizeLo             int64   `json:"DaySizeLo"`
	DaySizeHi             int64   `json:"DaySizeHi"`
	MonthSizeMB           float64 `json:"MonthSizeMB"`
	MonthSizeLo           int64   `json:"MonthSizeLo"`
	MonthSizeHi           int64   `json:"MonthSizeHi"`
	QuotaSizeMB           float64 `json:"QuotaSizeMB"`
	QuotaSizeLo           int64   `json:"QuotaSizeLo"`
	QuotaSizeHi           int64   `json:"QuotaSizeHi"`

	DownloadRate           float64 `json:"DownloadRate"`
	AverageDownloadRate    float64 `json:"AverageDownloadRate"`
	DownloadLimit          float64 `json:"DownloadLimit"`

	UpTimeSec      int64 `json:"UpTimeSec"`
	DownloadTimeSec int64 `json:"DownloadTimeSec"`
	ThreadCount    int64 `json:"ThreadCount"`
	PostJobCount   int64 `json:"PostJobCount"`
	UrlCount       int64 `json:"UrlCount"`
	QueueScriptCount int64 `json:"QueueScriptCount"`

	ServerStandBy    bool `json:"ServerStandBy"`
	DownloadPaused   bool `json:"DownloadPaused"`
	ServerPaused     bool `json:"ServerPaused"`
	PostPaused       bool `json:"PostPaused"`
	ScanPaused       bool `json:"ScanPaused"`
	QuotaReached     bool `json:"QuotaReached"`

	NewsServers []NewsServer `json:"NewsServers"`
}

type NewsServer struct {
	ID     int    `json:"ID"`
	Active bool   `json:"Active"`
}

func loHiBytes(lo, hi int64) float64 {
	return float64(hi)*4294967296 + float64(lo)
}

func (s *StatusResult) RemainingBytes() float64    { return loHiBytes(s.RemainingSizeLo, s.RemainingSizeHi) }
func (s *StatusResult) DownloadedBytes() float64   { return loHiBytes(s.DownloadedSizeLo, s.DownloadedSizeHi) }
func (s *StatusResult) ArticleCacheBytes() float64 { return loHiBytes(s.ArticleCacheLo, s.ArticleCacheHi) }
func (s *StatusResult) FreeDiskBytes() float64     { return loHiBytes(s.FreeDiskSpaceLo, s.FreeDiskSpaceHi) }
func (s *StatusResult) TotalDiskBytes() float64    { return loHiBytes(s.TotalDiskSpaceLo, s.TotalDiskSpaceHi) }
func (s *StatusResult) DayBytes() float64          { return loHiBytes(s.DaySizeLo, s.DaySizeHi) }
func (s *StatusResult) MonthBytes() float64        { return loHiBytes(s.MonthSizeLo, s.MonthSizeHi) }
