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

func TestMountPointPathname(t *testing.T) {
	pathname := MountPointPathname("teststring")
	expected := "YlOzkHHl34tQmPWSAtQUw3oX1qOKh171-MfYmwISsChpLT0gkM4Drh3mbIYvqKVh5X7Z63k1zmJzRPdCwJMdcg"
	if pathname != expected {
		t.Errorf("MountPointPathname was incorrect, got: %s, want: %s.", pathname, expected)
	}
}

func TestMountPointPathnameIdempotency(t *testing.T) {
	randVolumeName := randStringBytesMaskImprSrc(64)
	pathname1 := MountPointPathname(randVolumeName)
	pathname2 := MountPointPathname(randVolumeName)
	if pathname1 != pathname2 {
		t.Errorf("MountPointPathname was not identical with the same input, got: %s, want: %s.", pathname1, pathname2)
	}
}
