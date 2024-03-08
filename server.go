package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type RedisValue struct {
	value        string
	expiryMillis int64
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	database := make(map[string]RedisValue)

	handleClient := func(l net.Listener) {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		for {
			inputBuffer := make([]byte, 4096)
			bytesRead, readErr := conn.Read(inputBuffer)

			if readErr != nil || bytesRead <= 0 {
				return
			}

			command, parseErr := parseCommand(string(inputBuffer[:bytesRead]))

			if parseErr != nil {
				return
			}

			if command.Name == Echo {
				conn.Write([]byte(fmt.Sprintf("+%s\r\n", command.Args[0])))
			} else if command.Name == Ping {
				conn.Write([]byte("+PONG\r\n"))
			} else if command.Name == Set {
				expiryMillis := int64(0)
				if len(command.Args) == 4 {
					if strings.ToUpper(command.Args[2]) != "PX" {
						conn.Write([]byte("-ERR syntax error\r\n"))
						return
					}
					expiryMillis, _ = strconv.ParseInt(command.Args[3], 10, 64)
					expiryMillis += time.Now().UnixMilli()
				}
				redisValue := RedisValue{value: command.Args[1], expiryMillis: expiryMillis}
				database[command.Args[0]] = redisValue

				conn.Write([]byte("+OK\r\n"))
			} else if command.Name == Get {
				redisEntry, exists := database[command.Args[0]]

				if !exists {
					conn.Write([]byte("$-1\r\n"))
					return
				}

				if redisEntry.expiryMillis != 0 && redisEntry.expiryMillis < time.Now().UnixMilli() {
					delete(database, command.Args[0])
					conn.Write([]byte("$-1\r\n"))
					return
				}

				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(redisEntry.value), redisEntry.value)))
			}
		}
	}

	for {
		go handleClient(l)
	}
}

type CommandName string

const (
	Echo = "ECHO"
	Ping = "PING"
	Set  = "SET"
	Get  = "GET"
)

type Command struct {
	Name CommandName
	Args []string
}

func parseCommand(input string) (Command, error) {
	inputArgs, error := Parse([]byte(input))

	if error != nil {
		return Command{}, error
	}

	if len(inputArgs) == 0 {
		return Command{}, fmt.Errorf("empty input")
	}

	commandName := CommandName(strings.ToUpper(inputArgs[0]))

	if commandName == Echo {
		return Command{Name: Echo, Args: inputArgs[1:]}, nil
	} else if commandName == Ping {
		// PING command has no arguments
		return Command{Name: Ping, Args: nil}, nil
	} else if commandName == Set {
		if (len(inputArgs)-1)%2 != 0 {
			return Command{}, fmt.Errorf("invalid number of arguments")
		}
		return Command{Name: Set, Args: inputArgs[1:]}, nil
	} else if commandName == Get {
		if len(inputArgs) != 2 {
			return Command{}, fmt.Errorf("invalid number of arguments")
		}
		return Command{Name: Get, Args: inputArgs[1:]}, nil
	}

	return Command{}, fmt.Errorf("unknown command")
}
