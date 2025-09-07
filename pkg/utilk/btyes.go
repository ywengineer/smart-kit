package utilk

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"

	"github.com/cloudwego/netpoll"
)

func Int64ToBytes(odr binary.ByteOrder, i int64) []byte {
	var buf = make([]byte, 8)
	odr.PutUint64(buf, uint64(i))
	return buf
}

func Int32ToBytes(odr binary.ByteOrder, i int32) []byte {
	var buf = make([]byte, 4)
	odr.PutUint32(buf, uint32(i))
	return buf
}

func Int16ToBytes(odr binary.ByteOrder, i int16) []byte {
	var buf = make([]byte, 2)
	odr.PutUint16(buf, uint16(i))
	return buf
}

func BytesToInt64(odr binary.ByteOrder, buf []byte) int64 {
	return int64(odr.Uint64(buf))
}

func NewLinkBuffer(data []byte) *netpoll.LinkBuffer {
	lb := netpoll.NewLinkBuffer(len(data))
	_, _ = lb.WriteBinary(data)
	_ = lb.Flush()
	return lb
}

// Hash 类似 Java Objects.hash()，计算多个值的组合哈希
func Hash(values ...interface{}) uint64 {
	h := fnv.New64a()
	for _, v := range values {
		_, _ = fmt.Fprintf(h, "%v", v)
	}
	return h.Sum64()
}
