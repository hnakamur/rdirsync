package pb

import "time"

func ConvertTimeFromPB(t *UnixTime) time.Time {
	return time.Unix(t.Second, t.NanoSecond)
}

func ConvertTimeToPB(t time.Time) *UnixTime {
	return &UnixTime{
		Second:     t.Unix(),
		NanoSecond: t.UnixNano(),
	}
}
