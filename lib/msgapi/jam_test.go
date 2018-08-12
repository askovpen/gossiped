package msgapi

import(
  "testing"
)

func TestJamCrc32r(t *testing.T) {
  if crc32r("")!=0xffffffff {
    t.Errorf("got 0x%08x; want 0xffffffff", crc32r(""))
  }
  if crc32r("Alexander N. Skovpen")!=0x30222bd1 {
    t.Errorf("got 0x%08x; want 0x30222bd1", crc32r("Alexander N. Skovpen"))
  }
}
