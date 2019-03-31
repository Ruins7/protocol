package data

import (
	"math/big"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
)

// Handle an incoming string
func parseString(amount string, base *big.Float) *big.Int {
	//log.Dump("Parsing Amount", amount)
	if amount == "" {
		log.Error("Empty Amount String", "amount", amount)
		return nil
	}

	/*
		value := new(big.Float)

		_, err := fmt.Sscan(amount, value)
		if err != nil {
			log.Error("Invalid Float String", "err", err, "amount", amount)
			return nil
		}
		result := bfloat2bint(value, base)
	*/

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Error("Invalid Float String", "err", err, "amount", amount)
		return nil
	}

	result := float2bint(value, base)

	return result
}

func units2bint(amount int64, base *big.Float) *big.Int {
	value := new(big.Int).SetInt64(amount)
	return value
}

// Handle an incoming int (often used for comaparisons)
func int2bint(amount int64, base *big.Float) *big.Int {
	value := new(big.Float).SetInt64(amount)

	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	//log.Dump("int2bint", amount, result)
	return result
}

// Handle an incoming float
func float2bint(amount float64, base *big.Float) *big.Int {
	value := big.NewFloat(amount)

	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	return result
}

// Handle a big float to big int conversion
func bfloat2bint(value *big.Float, base *big.Float) *big.Int {
	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	return result
}

// Handle a big int to outgoing float
func bint2float(amount *big.Int, base *big.Float) float64 {
	value := new(big.Float).SetInt(amount)

	interim := value.Quo(value, base)
	result, _ := interim.Float64()

	return result
}
