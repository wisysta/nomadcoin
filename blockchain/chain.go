// 1. value receiver -> 객체를 value로 가져와서 사용하는 리시버
// 2. pointer receiver -> 객체를 reference(pointer)로 가져와서 사용하는 리시버
// 포인터란 값으로 메모리의 주소를 가지는 타입

package blockchain

import (
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis Block")
		})
	}

	return b
}
