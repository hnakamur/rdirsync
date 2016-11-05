package pb

import "time"

func ConvertTimeFromPB(t int64) time.Time {
	return time.Unix(t/1e9, t%1e9)
}

func ConvertTimeToPB(t time.Time) int64 {
	return t.UnixNano()
}
