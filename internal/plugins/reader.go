package plugins

import (
	"bytes"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/schneider001/sf-apis/go/sfgo"
)

const (
	// BuffSize represents the buffer size of the stream
	BuffSize = 16384
	// OOBuffSize represents the OO buffer size of the stream
	OOBuffSize = 1024
)

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
		log.Fatal("listen error:", err)
	}

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
			log.Fatal("accept error:", err)
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
		}

		conn.Close()
	}
}
