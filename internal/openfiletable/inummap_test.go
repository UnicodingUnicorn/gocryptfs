package openfiletable

import (
	"sync"
	"testing"
)

func TestTranslate(t *testing.T) {
	const baseDev = 12345
	m := NewInumMap(baseDev)

	q := QIno{Dev: baseDev, Ino: 1}
	out := m.Translate(q)
	if out != 1 {
		t.Errorf("expected 1, got %d", out)
	}
	q.Ino = inumTranslateBase
	out = m.Translate(q)
	if out < inumTranslateBase {
		t.Errorf("got %d", out)
	}
	out2 := m.Translate(q)
	if out2 != out {
		t.Errorf("unstable mapping: %d %d", out2, out)
	}
}

func TestTranslateStress(t *testing.T) {
	const baseDev = 12345
	m := NewInumMap(baseDev)
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		q := QIno{Dev: baseDev}
		for i := uint64(1); i <= 10000; i++ {
			q.Ino = i
			out := m.Translate(q)
			if out != i {
				t.Fail()
			}
		}
		wg.Done()
	}()
	go func() {
		q := QIno{Dev: baseDev}
		for i := uint64(1); i <= 10000; i++ {
			q.Ino = inumTranslateBase + i
			out := m.Translate(q)
			if out < inumTranslateBase {
				t.Fail()
			}
		}
		wg.Done()
	}()
	go func() {
		q := QIno{Dev: 9999999}
		for i := uint64(1); i <= 10000; i++ {
			q.Ino = i
			out := m.Translate(q)
			if out < inumTranslateBase {
				t.Fail()
			}
		}
		wg.Done()
	}()
	go func() {
		q := QIno{Dev: 4444444}
		for i := uint64(1); i <= 10000; i++ {
			q.Ino = i
			out := m.Translate(q)
			if out < inumTranslateBase {
				t.Fail()
			}
		}
		wg.Done()
	}()
	wg.Wait()
	if m.Count() != 30000 {
		t.Fail()
	}
}
