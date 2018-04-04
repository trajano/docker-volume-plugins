package main

import (
	"math/rand"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func TestMountPointPath(t *testing.T) {
	pathname := MountPointPath("teststring")
	expected := "/volumes/YlOzkHHl34tQmPWSAtQUw3oX1qOKh171-MfYmwISsChpLT0gkM4Drh3mbIYvqKVh5X7Z63k1zmJzRPdCwJMdcg"
	if pathname != expected {
		t.Errorf("MountPointPath was incorrect, got: %s, want: %s.", pathname, expected)
	}
}

func TestMountPointPathIdempotency(t *testing.T) {
	randVolumeName := randStringBytesMaskImprSrc(64)
	pathname1 := MountPointPath(randVolumeName)
	pathname2 := MountPointPath(randVolumeName)
	if pathname1 != pathname2 {
		t.Errorf("MountPointPath was not identical with the same input, got: %s, want: %s.", pathname1, pathname2)
	}
}
