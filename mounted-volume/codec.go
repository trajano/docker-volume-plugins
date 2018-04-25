package mountedvolume

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
)

const (
	volumeBucket = "volumes"
)

func (p *mountedVolumeInfo) gobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecode(data []byte) (*mountedVolumeInfo, error) {
	var p *mountedVolumeInfo
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Driver) storeVolumeInfo(tx *bolt.Tx, volumeName string, info *mountedVolumeInfo) error {
	bucket := tx.Bucket([]byte(volumeBucket))
	b, err := info.gobEncode()
	if err != nil {
		return err
	}
	return bucket.Put([]byte(volumeName), b)
}

func (p *Driver) getVolumeInfo(tx *bolt.Tx, volumeName string) (*mountedVolumeInfo, bool, error) {
	bucket := tx.Bucket([]byte(volumeBucket))
	v := bucket.Get([]byte(volumeName))
	if v == nil {
		return nil, false, nil
	}
	info, err := gobDecode(v)
	return info, true, err
}

func (p *Driver) getVolumeMap(tx *bolt.Tx) (map[string]mountedVolumeInfo, error) {
	bucket := tx.Bucket([]byte(volumeBucket))
	ret := make(map[string]mountedVolumeInfo)
	err := bucket.ForEach(func(k, v []byte) error {
		info, err := gobDecode(v)
		if err != nil {
			return err
		}
		ret[string(k)] = *info
		return nil
	})
	return ret, err
}

func (p *Driver) removeVolumeInfo(tx *bolt.Tx, volumeName string) error {
	bucket := tx.Bucket([]byte(volumeBucket))
	return bucket.Delete([]byte(volumeName))
}
