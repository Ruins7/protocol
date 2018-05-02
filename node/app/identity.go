/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains

	TODO: Need to pick a system key for identities. Is a hash of pubkey reasonable?
*/
package app

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data"
	"golang.org/x/crypto/ripemd160"
)

// Aliases to hide some of the basic underlying types.

type Address = data.Bytes // OneLedger address, like Tendermint the hash of the associated PubKey

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

type Signature = crypto.Signature

// enum for type
type IdentityType int

const (
	ONELEDGER IdentityType = iota
	BITCOIN
	ETHEREUM
)

type IdentityKey []byte

// Polymorphism
type Identity interface {
	AddPrivateKey(PrivateKey)
	Name() string
}

type IdentityBase struct {
	Type IdentityType

	Name string // TODO: Not sure this is normalized?

	Key       IdentityKey
	PublicKey PublicKey
}

// Information we need about our own fullnode identities
type IdentityOneLedger struct {
	IdentityBase

	PrivateKey PrivateKey

	NodeId      string
	ExternalIds []Identity
}

// Information we need for a Bitcoin account
type IdentityBitcoin struct {
	IdentityBase

	PrivateKey PrivateKey
}

// Information we need for an Ethereum account
type IdentityEthereum struct {
	IdentityBase

	PrivateKey PrivateKey
}

// Hash the public key to get a unqiue hash that can act as a key
func NewIdentityKey(key PublicKey) IdentityKey {
	hasher := ripemd160.New()

	bytes, err := key.MarshalJSON()
	if err != nil {
		panic("Unable to Marshal the key into bytes")
	}

	hasher.Write(bytes)

	return hasher.Sum(nil)
}

func NewIdentity(newType IdentityType, name string, Key PublicKey) Identity {
	switch newType {

	case ONELEDGER:
		return &IdentityOneLedger{}

	case BITCOIN:
		return &IdentityBitcoin{}

	case ETHEREUM:
		return &IdentityEthereum{}

	default:
		panic("Unknown Type")
	}
}

// TODO: really should be part of the enum, as a map...
func FindIdentityType(typeName string) (IdentityType, Error) {
	switch typeName {
	case "OneLedger":
		return ONELEDGER, 0

	case "Ethereum":
		return ETHEREUM, 0

	case "Bitcoin":
		return BITCOIN, 0
	}
	return 0, 42
}

func FindIdentity(name string) (Identity, Error) {
	// TODO: Lookup the identity in the node's database
	return &IdentityOneLedger{IdentityBase: IdentityBase{Name: name}}, 0
}

func (identity *IdentityOneLedger) AddPublicKey(key PublicKey) {
	identity.PublicKey = key
}

func (identity *IdentityOneLedger) AddPrivateKey(key PrivateKey) {
	identity.PrivateKey = key
}

func (identity *IdentityOneLedger) Name() string {
	return identity.IdentityBase.Name
}

func (identity *IdentityBitcoin) AddPublicKey(key PublicKey) {
	identity.PublicKey = key
}

func (identity *IdentityBitcoin) AddPrivateKey(key PrivateKey) {
	identity.PrivateKey = key
}

func (identity *IdentityBitcoin) Name() string {
	return identity.IdentityBase.Name
}

func (identity *IdentityEthereum) AddPublicKey(key PublicKey) {
	identity.PublicKey = key
}

func (identity *IdentityEthereum) AddPrivateKey(key PrivateKey) {
	identity.PrivateKey = key
}

func (identity *IdentityEthereum) Name() string {
	return identity.IdentityBase.Name
}