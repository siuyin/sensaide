package msg

import (
	"fmt"
	"testing"
)

func TestListTopics(t *testing.T) {
	ts := ListTopics()
	if len(ts) == 0 {
		t.Error("there should be topics")
	}
	for _, tp := range ts {
		fmt.Println(tp.ID())
	}
}
