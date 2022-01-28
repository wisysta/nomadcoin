package blockchain

import (
	"errors"
	"time"

	"github.io/wisysta/nomadcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

//TX
//	TxIn[$5(you)]
//	TxOut[$3(store)/$2(you)]
//	돈을 세는 방법이란 결과적으로 거래 출력값을 세어 보는것

type Tx struct {
	ID        string   `json:"id"`
	TimeStamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func isOnMempool(UTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == UTxOut.TxID && input.Index == UTxOut.Index {
				exists = true
				break Outer
			}
		}
	}

	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}

	txOuts := []*TxOut{
		{address, minerReward},
	}

	tx := Tx{
		ID:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enough money")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutByAddress(from, Blockchain())

	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.getId()

	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("nico", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) txToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("nico")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
