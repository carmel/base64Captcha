// Copyright 2017 Eric Zhou. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var (
	GCLimitNumber = 10240
	Expiration    = 10 * time.Minute
)

func TestSetGet(t *testing.T) {
	s := NewMemoryStore(GCLimitNumber, Expiration)
	id := "captcha id"
	d := "random-string"
	s.Set(id, d)
	d2 := s.Get(id, false)
	if d2 != d {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestGetClear(t *testing.T) {
	s := NewMemoryStore(GCLimitNumber, Expiration)
	id := "captcha id"
	d := "932839jfffjkdss"
	s.Set(id, d)
	d2 := s.Get(id, true)
	if d != d2 {
		t.Errorf("saved %v, getDigitsClear returned got %v", d, d2)
	}
	d2 = s.Get(id, false)
	if d2 != "" {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}

func TestCollect(t *testing.T) {
	// TODO(dchest): can't test automatic collection when saving, because
	// it's currently launched in a different goroutine.
	s := NewMemoryStore(10, -1)
	// create 10 ids
	ids := make([]string, 10)
	d := "fdjsij892jfi392j2"
	for i := range ids {
		ids[i] = fmt.Sprintf("%d", rand.Int63())
		s.Set(ids[i], d)
	}
	s.(*memoryStore).collect()
	// Must be already collected
	nc := 0
	for i := range ids {
		d2 := s.Get(ids[i], false)
		if d2 != "" {
			t.Errorf("%d: not collected", i)
			nc++
		}
	}
	if nc > 0 {
		t.Errorf("= not collected %d out of %d captchas", nc, len(ids))
	}
}

func BenchmarkSetCollect(b *testing.B) {
	b.StopTimer()
	d := "fdskfew9832232r"
	s := NewMemoryStore(9999, -1)
	ids := make([]string, 1000)
	for i := range ids {
		ids[i] = fmt.Sprintf("%d", rand.Int63())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			s.Set(ids[j], d)
		}
		s.(*memoryStore).collect()
	}
}

func TestMemoryStore_SetGoCollect(t *testing.T) {
	s := NewMemoryStore(10, -1)
	for i := 0; i <= 100; i++ {
		s.Set(fmt.Sprint(i), fmt.Sprint(i))
	}
}

func TestMemoryStore_CollectNotExpire(t *testing.T) {
	s := NewMemoryStore(10, time.Hour)
	for i := 0; i < 50; i++ {
		s.Set(fmt.Sprint(i), fmt.Sprint(i))
	}

	// let background goroutine to go
	time.Sleep(time.Second)

	assert.Equal(t, "0", s.Get("0", false))
}
