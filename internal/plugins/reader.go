package plugins

import (
	"bytes"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/schneider001/sf-apis/go/sfgo"
	"golang.org/x/sys/unix"
)

const (
	BuffSize   = 16384
	OOBuffSize = 1024
)

func checkSocketBuffer(conn *net.UnixConn) {
	rmemMax, err := getRmemMax()
	if err != nil {
		log.Fatal("Failed to get rmem_max:", err)
	}

	file, err := conn.File()
	if err != nil {
		log.Println("Failed to get file descriptor:", err)
		return
	}
	defer file.Close()

	const FIONREAD = 0x541B
	availableBytes, err := unix.IoctlGetInt(int(file.Fd()), FIONREAD)
	if err != nil {
		log.Println("Failed to get available bytes for reading from socket buffer:", err)
		return
	}

	log.Printf("Available bytes for reading from socket buffer: %d\n", availableBytes)
	log.Printf("rmem_max value: %d\n", rmemMax)

	if availableBytes >= rmemMax/2 {
		log.Fatalf("Available bytes for reading (%d) reached half of rmem_max (%d). Stopping program.", availableBytes, rmemMax)
	}
}

func getRmemMax() (int, error) {
	data, err := os.ReadFile("/proc/sys/net/core/rmem_max")
	if err != nil {
		return 0, err
	}

	rmemMaxStr := string(data)
	rmemMax, err := strconv.Atoi(strings.TrimSpace(rmemMaxStr))
	if err != nil {
		return 0, err
	}

	return rmemMax, nil
}

func mustSocket(socketPath string) {
	if _, err := os.Stat(socketPath); !errors.Is(err, os.ErrNotExist) {
		log.Println("Socket already exists")
		err = os.Remove(socketPath)
		if err != nil {
			log.Fatal("Socket remove: ", err)
		}
	} else {
		dir := filepath.Dir(socketPath)
		if err := os.MkdirAll(dir, 0600); err != nil {
			log.Fatal("Unable to create directory: ", err)
		}
	}
}

func Reader(socketPath string, records chan<- *sfgo.SysFlow) {
	mustSocket(socketPath)

	l, err := net.ListenUnix("unixpacket", &net.UnixAddr{Name: socketPath, Net: "unix"})
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	defer l.Close()

	sFlow := sfgo.NewSysFlow()
	deser, err := compiler.CompileSchemaBytes([]byte(sFlow.Schema()), []byte(sFlow.Schema()))
	if err != nil {
		log.Fatal("Compilation error: ", err)
	}

	for {
		health := false
		buf := make([]byte, BuffSize)
		oobuf := make([]byte, OOBuffSize)
		reader := bytes.NewReader(buf)

		conn, err := l.AcceptUnix()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		for {
			sFlow = sfgo.NewSysFlow()
			_, _, _, _, err := conn.ReadMsgUnix(buf[:], oobuf[:])
			if err != nil {
				log.Println("Read error: ", err)
				break
			}
			reader.Reset(buf)
			err = vm.Eval(reader, deser, sFlow)
			if err != nil {
				log.Println("Deserialization error: ", err)
				break
			}

			if !health {
				log.Println("Successfully read first record from input stream")
				health = true
			}

			records <- sFlow

			//checkSocketBuffer(conn)
		}

		conn.Close()
	}
}
