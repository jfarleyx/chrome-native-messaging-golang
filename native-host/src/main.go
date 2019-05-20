/*
* @Author: J. Farley
* @Date: 2019-05-19
* @Description: Basic chrome native messaging host example.
 */
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"unsafe"
)

// constants for Logger
var (
	// Trace logs general information messages.
	Trace *log.Logger
	// Error logs error messages.
	Error *log.Logger
)

// used to detect native byte order
var nativeEndian binary.ByteOrder

// Init initializes logger and determines native byte order.
func Init(traceHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// determine native byte order
	var one int16 = 1
	b := (*byte)(unsafe.Pointer(&one))
	if *b == 0 {
		nativeEndian = binary.BigEndian
	} else {
		nativeEndian = binary.LittleEndian
	}
}

// IncomingMessage represents a message sent to the native host.
type IncomingMessage struct {
	Query string `json:"query"`
}

// OutgoingMessage respresents a response to an incoming message query.
type OutgoingMessage struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}

func main() {
	file, err := os.OpenFile("chrome-native-host-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		Init(os.Stdout, os.Stderr)
		Error.Printf("Unable to create and/or open log file. Will log to Stdout and Stderr. Error: %v", err)
	} else {
		Init(file, file)
	}
	// ensure we close the log file when we're done
	defer file.Close()

	Trace.Print("Chrome native messaging host started.")
	Trace.Printf("Native byte order: %v", nativeEndian)
	read()
	Trace.Print("Chrome native messaging host exited.")
}

// read Creates a new buffered I/O reader and reads messages from Stdin.
func read() {
	v := bufio.NewReader(os.Stdin)
	// adjust buffer size to accommodate your json payload size limits; default is 4096
	s := bufio.NewReaderSize(v, 8192)
	Trace.Printf("IO buffer reader created with buffer size of %v.", s.Size())

	// we're going to read the first 4 bytes to get the message length
	lengthBytes := make([]byte, 4)
	lengthNum := int(0)

	for b, _ := s.Read(lengthBytes); b > 0; b, _ = s.Read(lengthBytes) {
		// get message length integer value
		lengthNum = readMessageLength(lengthBytes)
		Trace.Printf("Message total length: %v", lengthNum)

		// now read the content of the message; if content exceeds size of buffer in reader, this does not work
		content := make([]byte, lengthNum)
		_, err := s.Read(content)
		if err != nil && err != io.EOF {
			Error.Fatal(err)
		}

		// message has been read in full, now process it
		parseMessage(content)
	}

	Trace.Print("Stdin closed.")
}

// readMessageLength reads and returns the message length value in native byte order.
func readMessageLength(msg []byte) int {
	var length uint32
	buf := bytes.NewBuffer(msg)
	err := binary.Read(buf, nativeEndian, &length)
	if err != nil {
		Error.Printf("Unable to read bytes representing message length: %v", err)
	}
	return int(length)
}

// parseMessage parses incoming message and routes to appropriate process handlers.
func parseMessage(msg []byte) {
	iMsg := decodeMessage(msg)
	Trace.Printf("Message received: %s", msg)

	// start building outgoing json message
	oMsg := OutgoingMessage{
		Query: iMsg.Query,
	}

	switch iMsg.Query {
	case "ping":
		oMsg.Response = "pong"
	case "hello":
		oMsg.Response = "goodbye"
	default:
		oMsg.Response = "42"
	}

	send(oMsg)
}

// decodeMessage unmarshals incoming json request and returns query value.
func decodeMessage(msg []byte) IncomingMessage {
	var iMsg IncomingMessage
	err := json.Unmarshal(msg, &iMsg)
	if err != nil {
		Error.Printf("Unable to unmarshal json to struct: %v", err)
	}
	return iMsg
}

// send sends an OutgoingMessage to os.Stdout.
func send(msg OutgoingMessage) {
	byteMsg := dataToBytes(msg)
	writeMessageLength(byteMsg)

	var msgBuf bytes.Buffer
	_, err := msgBuf.Write(byteMsg)
	if err != nil {
		Error.Printf("Unable to write message length to message buffer: %v", err)
	}

	_, err = msgBuf.WriteTo(os.Stdout)
	if err != nil {
		Error.Printf("Unable to write message buffer to Stdout: %v", err)
	}
}

// dataToBytes marshals OutgoingMessage struct to slice of bytes
func dataToBytes(msg OutgoingMessage) []byte {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		Error.Printf("Unable to marshal OutgoingMessage struct to slice of bytes: %v", err)
	}
	return byteMsg
}

// writeMessageLength determines length of message and writes it to os.Stdout.
func writeMessageLength(msg []byte) {
	err := binary.Write(os.Stdout, binary.LittleEndian, uint32(len(msg)))
	if err != nil {
		Error.Printf("Unable to write message length to Stdout: %v", err)
	}
}
