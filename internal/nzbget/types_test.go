package nzbget

import "testing"

func TestLoHiBytes(t *testing.T) {
	// 1 * 2^32 + 0 = 4294967296
	got := loHiBytes(0, 1)
	if got != 4294967296 {
		t.Errorf("loHiBytes(0,1) = %v, want 4294967296", got)
	}
	// 0 + 500 = 500
	got = loHiBytes(500, 0)
	if got != 500 {
		t.Errorf("loHiBytes(500,0) = %v, want 500", got)
	}
}

func TestStatusResultMethods(t *testing.T) {
	s := &StatusResult{
		RemainingSizeLo:  100,
		RemainingSizeHi:  0,
		DownloadedSizeLo: 200,
		DownloadedSizeHi: 1,
	}
	if got := s.RemainingBytes(); got != 100 {
		t.Errorf("RemainingBytes = %v, want 100", got)
	}
	if got := s.DownloadedBytes(); got != 4294967496 {
		t.Errorf("DownloadedBytes = %v, want 4294967496", got)
	}
}
