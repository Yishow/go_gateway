package desktop

import (
	"errors"
	"testing"
	"time"
)

func TestBrowserTaskIsNonblockingAndOneInFlight(t *testing.T) {
	var task browserTask
	entered := make(chan string, 1)
	release := make(chan struct{})
	reported := make(chan struct{}, 1)
	opener := func(url string) error { entered <- url; <-release; return errors.New("untrusted browser detail") }
	if !task.start(opener, "fixed-url", func() { reported <- struct{}{} }) {
		t.Fatal("first opener rejected")
	}
	select {
	case url := <-entered:
		if url != "fixed-url" {
			t.Fatal(url)
		}
	case <-time.After(time.Second):
		t.Fatal("opener did not start")
	}
	if task.start(opener, "other-url", nil) {
		t.Fatal("unbounded concurrent opener accepted")
	}
	close(release)
	select {
	case <-reported:
	case <-time.After(time.Second):
		t.Fatal("fixed error notification missing")
	}
	deadline := time.Now().Add(time.Second)
	for task.running.Load() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if task.running.Load() {
		t.Fatal("opener stayed busy after completion")
	}
}
