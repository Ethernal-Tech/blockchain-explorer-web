package utils

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/params"
)

func ToUint64(str string) uint64 {
	if len(str) == 0 {
		return 0
	}

	var res uint64
	var err error

	if str[0:2] == "0x" {
		if len(str) <= 2 {
			return 0
		}

		res, err = strconv.ParseUint(str[2:], 16, 64)
		if err != nil {
			//TODO: error handling
			return 0
		}
	} else {
		res, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			//TODO: error handling
			return 0
		}
	}
	return res
}

func ToBigInt(str string) *big.Int {
	if len(str) == 0 {
		return big.NewInt(0)
	}

	res := big.NewInt(0)
	var err bool

	if str[0:2] == "0x" {
		if len(str) <= 2 {
			return big.NewInt(0)
		}

		i := new(big.Int)
		res, err = i.SetString(str[2:], 16)
		if !err {
			//TODO: error handling
			return big.NewInt(0)
		}
	} else {
		i := new(big.Int)
		res, err = i.SetString(str, 10)
		if !err {
			//TODO: error handling
			return big.NewInt(0)
		}
	}
	return res
}

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}
