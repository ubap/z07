package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"math/big"
)

// The constants you provided, represented as strings for parsing.
const (
	pStr  = "14299623962416399520070177382898895550795403345466153217470516082934737582776038882967213386204600674145392845853859217990626450972452084065728686565928113"
	qStr  = "7630979195970404721891201847792002125535401292779123937207447574596692788513647179235335529307251350570728407373705564708871762033017096809910315212884101"
	dpStr = "11141736698610418925078406669215087697114858422461871124661098818361832856659225315773346115219673296375487744032858798960485665997181641221483584094519937"
	dqStr = "4886309137722172729208909250386672706991365415741885286554321031904881408516947737562153523770981322408725111241551398797744838697461929408240938369297973"
	eInt  = 65537 // Common public exponent
)

func main() {
	// --- 1. Parse all the input strings into big.Int objects ---
	p, ok := new(big.Int).SetString(pStr, 10)
	if !ok {
		log.Fatal("Failed to parse p")
	}

	q, ok := new(big.Int).SetString(qStr, 10)
	if !ok {
		log.Fatal("Failed to parse q")
	}

	dp, ok := new(big.Int).SetString(dpStr, 10)
	if !ok {
		log.Fatal("Failed to parse dp")
	}

	dq, ok := new(big.Int).SetString(dqStr, 10)
	if !ok {
		log.Fatal("Failed to parse dq")
	}

	e := big.NewInt(eInt)

	// --- 2. Calculate the missing RSA components ---
	n := new(big.Int).Mul(p, q)
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
	phi := new(big.Int).Mul(pMinus1, qMinus1)
	d := new(big.Int).ModInverse(e, phi)
	qInv := new(big.Int).ModInverse(q, p)

	// --- 3. Assemble the rsa.PrivateKey struct ---
	privateKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: int(eInt),
		},
		D:      d,
		Primes: []*big.Int{p, q},
		Precomputed: rsa.PrecomputedValues{
			Dp:   dp,
			Dq:   dq,
			Qinv: qInv,
		},
	}

	// --- 4. Validate the key to ensure correctness ---
	if err := privateKey.Validate(); err != nil {
		log.Fatalf("Key validation failed: %v", err)
	}
	fmt.Println("✅ RSA Private Key components successfully validated!")

	// --- 5. Print the fully calculated components as usable Go constants ---
	fmt.Println("\n// You can copy and paste the following constants into your Go code.")

	fmt.Println("\n// --- RSA Public Key Components ---")
	fmt.Printf("const N = \"%s\"\n", privateKey.N)
	fmt.Printf("const E = %d\n", privateKey.E)

	fmt.Println("\n// --- RSA Private Key Components (Primary) ---")
	fmt.Printf("const D = \"%s\"\n", privateKey.D)

	fmt.Println("\n// --- RSA Private Key Components (Primes) ---")
	fmt.Printf("const P = \"%s\"\n", privateKey.Primes[0])
	fmt.Printf("const Q = \"%s\"\n", privateKey.Primes[1])

	fmt.Println("\n// --- RSA Private Key Components (CRT for Optimization) ---")
	fmt.Printf("const Dp = \"%s\"\n", privateKey.Precomputed.Dp)
	fmt.Printf("const Dq = \"%s\"\n", privateKey.Precomputed.Dq)
	fmt.Printf("const Qinv = \"%s\"\n", privateKey.Precomputed.Qinv)
}
