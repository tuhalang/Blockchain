package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Balances map[Account]uint
	txMemPool []Tx

	dbFile *os.File
}

func NewStateFromDisk() (*State, error){

	// get current working directory
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genFilePath)

	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)

	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{
		balances,
		make([]Tx, 0),
		f,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMemPool = append(s.txMemPool, tx)

	return nil
}

func (s *State) Persist() error {
	mempool := make([]Tx, len(s.txMemPool))
	copy(mempool, s.txMemPool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])

		if err != nil {
			return err
		}

		if _, err := s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		s.txMemPool = append(s.txMemPool[:i], s.txMemPool[i+1:]...)
	}

	return nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
