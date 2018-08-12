package msgapi

import(
  "testing"
)

func TestSquishBufHash32(t *testing.T) {
  if bufHash32("")!=0x0 {
    t.Errorf("got 0x%08x; want 0x00000000", bufHash32(""))
  }
  if bufHash32("Alexander N. Skovpen")!=0x00efd7be {
    t.Errorf("got 0x%08x; want 0x00efd7be", bufHash32("Alexander N. Skovpen"))
  }
  if bufHash32("Serg Ageev")!=0x0967dfc6 {
    t.Errorf("got 0x%08x; want 0x0967dfc6", bufHash32("Serg Ageev"))
  }
}