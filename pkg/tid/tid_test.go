package tid

import (
	"testing"
	"time"
)

var tidTests = []struct {
	TID     string
	Time    string
	ClockID int
}{
	{"2222222222222", "1970-01-01T00:00:00.000000-00:00", 0},
	{"jzzzzzzzzzzzz", "2540-11-07T23:35:09.481983-00:00", 1023},
	{"3mcvuiwkquc22", "2026-01-21T04:47:07.357000-00:00", 0},
	{"3mcvuiwkquc27", "2026-01-21T04:47:07.357000-00:00", 5},
	{"3mcvuiwkquczz", "2026-01-21T04:47:07.357000-00:00", 1023},
}

func TestTIDEncoding(t *testing.T) {
	for _, test := range tidTests {
		t0, err := time.Parse(time.RFC3339, test.Time)
		if err != nil {
			t.Error("invalid test time")
		}
		tid := TID(t0.UTC(), uint16(test.ClockID))
		if tid != test.TID {
			t.Errorf("incorrect encoding: expected %s got %s", test.TID, tid)
		}
	}
}

func TestTIDDecoding(t *testing.T) {
	for _, test := range tidTests {
		t0, err := time.Parse(time.RFC3339, test.Time)
		if err != nil {
			t.Error("invalid test time")
		}
		t0 = t0.UTC()
		t1, clockid, err := ParseTID(test.TID)
		if err != nil {
			t.Errorf("failed to parse %s", test.TID)
		}
		if t1 != t0 {
			t.Errorf("error parsing %s: expected %s got %s", test.TID, t0, t1)
		}
		if clockid != uint16(test.ClockID) {
			t.Errorf("error parsing %s: expected %d got %d", test.TID, test.ClockID, clockid)
		}
		t2, err := FromString(t0.Format(ATProtoTimeFormat), uint16(test.ClockID))
		if err != nil {
			t.Errorf("failed to parse %s", test.Time)
		}
		if t2 != test.TID {
			t.Errorf("error generating tid: expected %s got %s", test.TID, t2)
		}
	}
}
