package collector

import (
	"errors"
	"strings"
	"testing"

	"github.com/genma/nzbget-metrics/internal/nzbget"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type fakeClient struct {
	status  *nzbget.StatusResult
	version string
	queued  int
	err     error
}

func (f *fakeClient) Status() (*nzbget.StatusResult, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.status, nil
}

func (f *fakeClient) Version() (string, error) { return f.version, nil }
func (f *fakeClient) QueuedCount() (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.queued, nil
}

func TestCollectUp(t *testing.T) {
	c := New(&fakeClient{
		status: &nzbget.StatusResult{
			DownloadRate:    512000,
			RemainingSizeLo: 1000,
			NewsServers:     []nzbget.NewsServer{{ID: 1, Active: true}},
		},
		version: "21.1",
		queued:  3,
	})

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	problems, err := testutil.GatherAndLint(reg)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) > 0 {
		t.Errorf("lint problems: %v", problems)
	}

	gathered, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}

	byName := make(map[string]float64)
	for _, mf := range gathered {
		for _, m := range mf.GetMetric() {
			byName[mf.GetName()] = m.GetGauge().GetValue()
		}
	}

	if byName["nzbget_up"] != 1 {
		t.Errorf("nzbget_up = %v, want 1", byName["nzbget_up"])
	}
	if byName["nzbget_queued"] != 3 {
		t.Errorf("nzbget_queued = %v, want 3", byName["nzbget_queued"])
	}
}

func TestCollectDown(t *testing.T) {
	c := New(&fakeClient{err: errors.New("connection refused")})

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	if err := testutil.GatherAndCompare(reg, strings.NewReader(`
# HELP nzbget_up 1 if the nzbget API is reachable.
# TYPE nzbget_up gauge
nzbget_up 0
`), "nzbget_up"); err != nil {
		t.Fatal(err)
	}
}
