package data

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type KeyValueDB struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewDB(hostPort string) (key_value_db *KeyValueDB) {

	var err error = nil
	key_value_db = new(KeyValueDB)

	//Here we will start a TCP connection to the server
	key_value_db.conn, err = net.Dial("tcp", hostPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	key_value_db.reader = bufio.NewReader(key_value_db.conn)
	key_value_db.writer = bufio.NewWriter(key_value_db.conn)

	return
}

func (kv *KeyValueDB) Get(key string) (results []string) {
	var q string = fmt.Sprintf("get %s\n", key)
	kv.conn.Write([]byte(q))
	result, _ := kv.reader.ReadString(0)
	results = strings.Split(result, "\n")
	return
}

func (kv *KeyValueDB) GetAll(key string) (result []string) {
	result = make([]string, 0)

	for i := 0; i < 32; i++ {
		var q string = fmt.Sprintf("%s:%d", key, i)
		var i_result []string = kv.Get(q)

		for i := 0; i < len(i_result)-1; i++ {

			if strings.TrimSpace(i_result[i]) == "" {
				continue
			}
			result = append(result, i_result[i])
		}
	}

	return
}

func (kv *KeyValueDB) Set(key string, value []byte) (result string) {
	var q string = fmt.Sprintf("set %s %s\n", key, value)

	kv.conn.Write([]byte(q))

	result, _ = kv.reader.ReadString('\n')
	result = strings.TrimSpace(result)

	return
}

func (kv *KeyValueDB) Write() (result string) {

	kv.conn.Write([]byte("write\n"))

	result, _ = kv.reader.ReadString('\n')
	result = result[0 : len(result)-1]

	return
}

func (kv *KeyValueDB) exit() (result string) {

	kv.conn.Write([]byte("exit\n"))

	result, _ = kv.reader.ReadString('\n')
	result = result[0 : len(result)-1]

	return
}

func (kv *KeyValueDB) Close() {
	fmt.Printf("%s", kv.Write())
	kv.exit()
	kv.conn.Close()
}
