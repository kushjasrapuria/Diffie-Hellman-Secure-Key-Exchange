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
	serverAddr := "localhost:8080"

	// Resolve server address
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	n := generatePrime()
	y := generatePrime()

	// Buffer
	buf := make([]byte, 1024)

	// Message to send
	message := []byte("---Secure Key Exchange---")
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		os.Exit(1)
	}

	// Receive response
	m, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(buf[:m]))

	// Send n
	rn := []byte(BIntToBytes(n))
	_, err = conn.Write(rn)
	if err != nil {
		fmt.Println("Error sending message:", err)
		os.Exit(1)
	}

	// Read g
	tg, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		os.Exit(1)
	}
	g := (BytesToBInt(buf[:tg]))

	// Calculation
	b := calculateSharedSecret(n, y, g)

	// Read a
	ta, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		os.Exit(1)
	}
	a := (BytesToBInt(buf[:ta]))

	// Send b
	rb := []byte(BIntToBytes(b))
	_, err = conn.Write(rb)
	if err != nil {
		fmt.Println("Error sending message:", err)
		os.Exit(1)
	}

	// Key Calculation
	k := calculateSharedSecret(n, y, a)
	fmt.Printf("Key : %d", k)
}
