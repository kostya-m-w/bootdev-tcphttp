package main
import(
	"net"
	"os"
	"fmt"
	"bufio"
)

func main() {
	udpaddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error on resolving address: %v\n", err)
	}
	
	conn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		fmt.Printf("Error on dialing UDP: %v\n", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading STDIN: %v", err)
		}
		_, err = conn.Write([]byte(line)) 
		if err != nil {
			fmt.Printf("Error sending throug UDP: %v", err)
		}
	}

}
