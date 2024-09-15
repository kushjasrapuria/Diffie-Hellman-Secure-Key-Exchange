package main

import (
	"fmt"
	"math/big"
	"crypto/rand"
	"net"
	"os"
)

func generatePrime() *big.Int {
	prime, _ := rand.Prime(rand.Reader, 256) // Generating a 256-bit prime number
	return prime
}

func calculateSharedSecret(prime, privateKey, publicKey *big.Int) *big.Int {
	sharedSecret := new(big.Int).Exp(publicKey, privateKey, prime)
	return sharedSecret
}

func BIntToBytes(bI *big.Int) []byte {
    return []byte(fmt.Sprintf("%x", bI))
}

func BytesToBInt(bts []byte) *big.Int {
    var bI big.Int
    fmt.Sscanf(string(bts), "%x", &bI)
    return &bI
}

func main() {
	// Resolve UDP address
	addr := "localhost:8080"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// Create a UDP connection
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	g := generatePrime()
	x := generatePrime()
	
	// Buffer to hold received data
	buf := make([]byte, 1024)

	for {
		// Read response
		m, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}
		fmt.Printf("%s\n", string(buf[:m]))

		// Message to send
		r := []byte("---Secure Key Exchange---")
		_, err = conn.WriteToUDP(r, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}

		// Read n
		bn, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}
		n := BytesToBInt(buf[:bn])

		// Send g
		rg := BIntToBytes(g)
		_, err = conn.WriteToUDP(rg, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}

		// Calculation
		a := calculateSharedSecret(n, x, g)

		// Send a
		ra := BIntToBytes(a)
		_, err = conn.WriteToUDP(ra, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}

		// Read b
		bb, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}
		b := BytesToBInt(buf[:bb])

		// Key Calculation
		k := calculateSharedSecret(n, x, b)
		fmt.Printf("Key : %d\n", k)

		break
	}
}
