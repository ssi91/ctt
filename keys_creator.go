package ctt

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
)

type AllocatedSimulatedBackend struct {
	backends.SimulatedBackend
	allocatedKeys *[]*ecdsa.PrivateKey
}

//type SmartContract struct {
//	Contract        *interface{}
//	ContractAddress common.Address
//	OwnerKey        *ecdsa.PrivateKey
//	SystemKey       *ecdsa.PrivateKey
//}

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

func SendEtherToAddress(ctx context.Context, backend *AllocatedSimulatedBackend, senderKey *ecdsa.PrivateKey, receiverAddress *common.Address, value *big.Int) (*types.Receipt, error) {
	nonce, err := backend.PendingNonceAt(ctx, crypto.PubkeyToAddress(senderKey.PublicKey))
	if err != nil {
		log.Printf("SendEtherToAddress error: %s\n", err)
		return nil, err
	}

	var gasLimit uint64 = 2100000

	gasPrice, err := backend.SuggestGasPrice(ctx)
	if err != nil {
		log.Printf("SendEtherToAddress gas price error: %s\n", err)
		return nil, err
	}

	var data []byte
	tx := types.NewTransaction(nonce, *receiverAddress, value, gasLimit, gasPrice, data)

	chainID := big.NewInt(1337)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), senderKey)
	if err != nil {
		log.Printf("SendEtherToAddress sign error: %s\n", err)
		return nil, err
	}

	err = backend.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Printf("SendEtherToAddress SendTransaction: %s\n", err)
		return nil, err
	}

	backend.Commit()
	//log.Printf("gas costed: %s", signedTx.Data())

	return backend.TransactionReceipt(ctx, signedTx.Hash())
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
