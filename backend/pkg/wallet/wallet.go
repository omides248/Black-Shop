package wallet

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/go-bip39"
)

type Service struct {
	masterKey *hdkeychain.ExtendedKey
	mnemonic  string
}

func NewService(mnemonic string) (*Service, error) {

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic phrase")
	}

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	return &Service{
		masterKey: masterKey,
		mnemonic:  mnemonic,
	}, nil
}

func (s *Service) DeriveAddress(path string) (string, error) {

	// 1. Parse the derivation path string (e.g., "m/44'/60'/0'/0/0")
	parts, err := parseDerivationPath(path)
	if err != nil {
		return "", err
	}

	// 2. Start with the master key and derive child keys sequentially
	currentKey := s.masterKey

}
