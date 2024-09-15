package main

import (
	"fmt"
	"math/big"
	"math/rand/v2"
	"math"
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

	g := getPrime(2, 15)	
	x := getPrime(2, 15)
	for {
		if g == x {
			x = getPrime(10, 20)
		} else {
			break
		}
	}
	
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
		n := BytesToInt(buf[:bn])

		// Send g
		rg := IntToBytes(g)
		_, err = conn.WriteToUDP(rg, addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}

		// Calculation
		a := (IntPow(g, x)) % n

		// Send a
		ra := IntToBytes(int(a))
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
		b := BytesToInt(buf[:bb])

		// Key Calculation
		k := (IntPow(b, x)) % n
		fmt.Printf("Key : %d\n", k)

		break
	}
}
