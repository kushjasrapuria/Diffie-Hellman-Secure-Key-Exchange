package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"net"
	"os"
)

func IsPrime(value int) bool {
    for i := 2; i <= int(math.Floor(float64(value) / 2)); i++ {
        if value%i == 0 {
            return false
        }
    }
    return true
}

func getPrime(minval, maxval int) int {
	x := rand.IntN(maxval-minval) + minval
	for {
		if IsPrime(x) == false {
			x = rand.IntN(maxval-minval) + minval
		} else {
			return x
		}
	}
}

func IntToBytes(i int) []byte{
    if i > 0 {
        return append(big.NewInt(int64(i)).Bytes(), byte(1))
    }
    return append(big.NewInt(int64(i)).Bytes(), byte(0))
}

func BytesToInt(b []byte) int{
    if b[len(b)-1]==0 {
        return -int(big.NewInt(0).SetBytes(b[:len(b)-1]).Int64())
    }
    return int(big.NewInt(0).SetBytes(b[:len(b)-1]).Int64())
}

func IntPow(n, m int) int {

    if m == 0 {
        return 1
    }

    if m == 1 {
        return n
    }

    result := n
    for i := 2; i <= m; i++ {
        result *= n
    }
    return result
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

	n := getPrime(2, 15)
	y := getPrime(2, 15)
	for {
		if n == y {
			y = getPrime(10, 20)
		} else {
			break
		}
	}

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
	rn := []byte(IntToBytes(n))
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
	g := (BytesToInt(buf[:tg]))

	// Calculation
	b := (IntPow(g, y)) % n

	// Read a
	ta, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		os.Exit(1)
	}
	a := (BytesToInt(buf[:ta]))

	// Send b
	rb := []byte(IntToBytes(int(b)))
	_, err = conn.Write(rb)
	if err != nil {
		fmt.Println("Error sending message:", err)
		os.Exit(1)
	}

	// Key Calculation
	k := (IntPow(a, y)) % n
	fmt.Printf("Key : %d", k)
}
