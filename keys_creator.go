package ctt

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
)

type AllocatedSimulatedBackend struct {
	backends.SimulatedBackend
	allocatedKeys *[]*ecdsa.PrivateKey
}

func NewAllocatedSimulatedBackend(amount uint) *AllocatedSimulatedBackend {
	alloc, keys := generateKeysAlloc(amount)

	var gasLimit uint64 = 9999999
	backend := backends.NewSimulatedBackend(*alloc, gasLimit)

	return &AllocatedSimulatedBackend{
		SimulatedBackend: *backend,
		allocatedKeys:    keys,
	}
}

func (b *AllocatedSimulatedBackend) Backend() backends.SimulatedBackend {
	return b.SimulatedBackend
}

func (b *AllocatedSimulatedBackend) Keys() *[]*ecdsa.PrivateKey {
	if b.allocatedKeys == nil {
		log.Fatal("Backend's Keys are not initialized")
	}
	return b.allocatedKeys
}

func generateKeysAlloc(amount uint) (*core.GenesisAlloc, *[]*ecdsa.PrivateKey) {
	alloc := make(core.GenesisAlloc)
	keys := make([]*ecdsa.PrivateKey, amount)
	balance := new(big.Int)
	balance.SetString("100000000000000000000", 10) // 100 eth in wei
	for i := uint(0); i < amount; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Generate key error: %s\n", err)
		}

		keys[i] = key

		alloc[crypto.PubkeyToAddress(key.PublicKey)] = core.GenesisAccount{
			Balance: balance,
		}
	}

	return &alloc, &keys
}
