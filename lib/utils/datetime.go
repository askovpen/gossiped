package utils

import (
  "time"
)

func F2T(d uint32) time.Time {
  return time.Date( int((d >> 9) & 127)+1980,
    time.Month(int(d>>5&15)),
    int(d&31),
    int(d>>27&31),
    int(d>>21&63),
    int(d>>16&31)*2,
    0,
    time.Local)
}
