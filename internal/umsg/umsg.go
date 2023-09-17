/*
umsg - UART messaging

# Once configured the user of this package can publish message message to other devices via UART

DEVTODO - describe how devices are connected with on uart use for input and the other uart used for output
  - describe how the devices are expected to make a loop so that a message sent is forwarded around the loop until returns to the sender
  - describe how users can subscribe to messages by listening to specific message queues (aka channels)
*/
package umsg

import (
	"log"
	"machine"
	"runtime"
	"strings"
	"time"
)

// Define message types
type MsgType string

const (
	TOKEN_HAT         byte   = 94  // ^
	TOKEN_ABOUT       byte   = 126 // ~
	LOOKBACK_SENDERID string = "Loopback"
)

const (
	MSG_FOO MsgType = "Foo"
	MSG_BAR MsgType = "Bar"
)

// ^Foo|some-sender|This is a foo message~
type FooMsg struct {
	Kind     MsgType
	SenderID string
	Name     string
}
type BarMsg struct {
	Kind     MsgType
	SenderID string
	Name     string
}

type MsgInterface interface {
	FooMsg | BarMsg
}

type UART interface {
	Configure(config machine.UARTConfig) error
	Buffered() int
	ReadByte() (byte, error)
	Write(data []byte) (n int, err error)
}

// Message Broker
type MsgBroker struct {
	// Sender ID set on each message
	// If a sender receives its own message it will not be forwarded on
	senderID string

	uartIn      UART
	uartInTxPin machine.Pin
	uartInRxPin machine.Pin

	uartOut      UART
	uartOutTxPin machine.Pin
	uartOutRxPin machine.Pin

	fooCh chan FooMsg
	barCh chan BarMsg
}

func NewBroker(
	senderID string,

	uartIn UART,
	uartInTxPin machine.Pin,
	uartInRxPin machine.Pin,

	uartOut UART,
	uartOutTxPin machine.Pin,
	uartOutRxPin machine.Pin,

	fooCh chan FooMsg,
	barCh chan BarMsg,

) MsgBroker {

	var mb MsgBroker

	mb.senderID = senderID

	if uartIn != nil {
		mb.uartIn = uartIn
		mb.uartInTxPin = uartInTxPin
		mb.uartInRxPin = uartInRxPin
	}

	if uartOut != nil {
		mb.uartOut = uartOut
		mb.uartOutTxPin = uartOutTxPin
		mb.uartOutRxPin = uartOutRxPin
	}

	if fooCh != nil {
		mb.fooCh = fooCh
	}

	if barCh != nil {
		mb.barCh = barCh
	}

	return mb

}

func (mb *MsgBroker) Configure() {

	// Output UART
	if mb.uartOut != nil {

		mb.uartOut.Configure(machine.UARTConfig{TX: mb.uartOutTxPin, RX: mb.uartOutRxPin})

	}

	// Input UART
	if mb.uartIn != nil {

		mb.uartIn.Configure(machine.UARTConfig{TX: mb.uartInTxPin, RX: mb.uartInRxPin})
		// Launch the listenRoutine to watch the input UART
		go mb.listenRoutine()

	}

}

// DEVTODO - make this generic
func (mb *MsgBroker) PublishFoo(foo FooMsg) {

	if foo.SenderID == "" {
		foo.SenderID = mb.senderID
	}

	msgStr := "^" + string(foo.Kind)
	msgStr = msgStr + "|" + foo.SenderID
	msgStr = msgStr + "|" + foo.Name + "~"

	mb.writeMsg(msgStr)

}

func (mb *MsgBroker) PublishBar(bar BarMsg) {

	if bar.SenderID == "" {
		bar.SenderID = mb.senderID
	}

	msgStr := "^" + string(bar.Kind)
	msgStr = msgStr + "|" + bar.SenderID
	msgStr = msgStr + "|" + bar.Name + "~"

	mb.writeMsg(msgStr)

}

func (mb *MsgBroker) writeMsg(msg string) {

	if mb.uartOut != nil {

		log.Printf("[writeMsg] - sending message: %v\n", msg)
		mb.uartOut.Write([]byte(msg))
		// Print a new line between messages for readability in the serial monitor
		mb.uartOut.Write([]byte("\n"))

	} else {

		log.Printf("[writeMsg] - message not sent, no output uart\n")

	}

}

func (mb *MsgBroker) dispatchMsgToChannel(msgParts []string) {

	switch msgParts[0] {

	case string(MSG_FOO):
		log.Printf("[dispatchMsgToChannel] - %v\n", MSG_FOO)
		msg := makeFoo(msgParts)
		log.Printf("[dispatchMsgToChannel] - msg: %v\n", msg)
		if mb.fooCh != nil {
			log.Printf("[dispatchMsgToChannel] - write to fooCh: %v\n", *msg)
			mb.fooCh <- *msg
		} else {
			log.Printf("[dispatchMsgToChannel] - send to bit bucket, no fooCh: %v\n")
		}
	case string(MSG_BAR):
		log.Printf("[dispatchMsgToChannel] - %v\n", MSG_BAR)
		msg := makeBar(msgParts)
		if mb.barCh != nil {
			mb.barCh <- *msg
		}

	default:
		log.Println("[dispatchMsgToChannel] - no match found")
	}

}

// DEVTODO - make this generic
func makeFoo(msgParts []string) *FooMsg {

	fooMsg := new(FooMsg)

	if len(msgParts) > 0 {
		fooMsg.Kind = MSG_FOO
	}
	if len(msgParts) > 1 {
		fooMsg.SenderID = msgParts[1]
	}
	if len(msgParts) > 2 {
		fooMsg.Name = msgParts[2]
	}

	return fooMsg
}

func makeBar(msgParts []string) *BarMsg {

	barMsg := new(BarMsg)

	if len(msgParts) > 0 {
		barMsg.Kind = MSG_BAR
	}
	if len(msgParts) > 1 {
		barMsg.SenderID = msgParts[1]
	}
	if len(msgParts) > 2 {
		barMsg.Name = msgParts[2]
	}

	return barMsg
}

/*
listenRoutine will monitor the input uart for messages and dispatch each message
to a specific channel based on the message type
*/
func (mb *MsgBroker) listenRoutine() {

	log.Println("[listenRoutine] - Start listenRoutine loop...")
	for {

		log.Println("[listenRoutine] - call readMsg()")
		msg, more := mb.readMsg()
		msgParts := strings.Split(string(msg), "|")

		log.Println("[listenRoutine] - Check if empty message")
		if len(msgParts) > 2 {

			// Get the message senderID
			
			// DEVTOD for now it is assumed that index 1 is sender id
			msgSenderID := msgParts[1]

			// Only dispatch messages from other senders
			if msgSenderID != mb.senderID {

				mb.dispatchMsgToChannel(msgParts)

				// Forward all messages with the exception of the loopback sender to prevent endless loop
				if mb.uartOut != nil && msgSenderID != LOOKBACK_SENDERID {
					// rewrap the message to start with ^ and end with ~
					msg = string(TOKEN_HAT) + msg + string(TOKEN_ABOUT)
					log.Printf("[listenRoutine] send message to output uart: %v\n", msg)
					mb.uartOut.Write([]byte(msg))
				} else {
					log.Printf("[listenRoutine] drop message because there is no output uart or this is a loopback sender: %v\n", msg)
				}

			} else {
				log.Printf("[listenRoutine] drop message because senderID same as broker senderID: %v\n", msg)
			}

		} else {
			log.Printf("[listenRoutine] - no message or not enough parts, msgParts: %v", msgParts)
		}

		
		// If there are no more messages in the buffer then wait before trying again
		// otherwise try again without delay
		if !more {						
			// DEVTODO - what is is a good delay time? I don't want to run down the battery
			runtime.Gosched()
			time.Sleep(time.Millisecond * 2000)
		}



	}
}

/*
readMsg will read the input buffer looking for a message

Given:

	this-is-junk^Foo|some-sender|This is a foo message~^Bar|some-sender|This is a bar message~more-junk

The following string is returned:

	Foo|some-sender|This is a foo message

The next time readMsg() is called this is returned:

	Bar|some-sender|This is a bar message
*/
func (mb *MsgBroker) readMsg() (msg string, more bool) {

	// used to hold message read from input UART
	message := make([]byte, 0)

	// Seek receive buffer to start of next message
	// if no message is found then get out
	log.Println("[readMsg] - calling seekStartOfMessage()")
	if !mb.seekStartOfMessage() {
		log.Println("[readMsg] - did not find start of message")
		return "", false
	}

	//
	// Start read message loop
	//
	log.Println("[readMsg] - start read loop")
	for {

		// Read from buffer
		data, ok := mb.uartIn.ReadByte()

		if data == TOKEN_ABOUT {
			log.Printf("[readMsg] - break out of read loop because we hit the end of message; data: [%v]", data)
			break
		} 

		if ok != nil  {
			log.Printf("[readMsg] - end of buffer hit before we found the end of message, pause then read more...")
			runtime.Gosched()
			time.Sleep(time.Millisecond * 100)
		} else {
			message = append(message, data)
		}
	}

	// Set return values
	if len(message) > 0 {
		log.Printf("[readMsg] - return this message:  %v\n", string(message))
		return string(message), mb.uartIn.Buffered() > 0
	} else {
		log.Printf("[readMsg] - return empty message\n")
		return "", mb.uartIn.Buffered() > 0
	}

}

func (mb *MsgBroker) seekStartOfMessage() (isFound bool) {

	log.Printf("[seekStartOfMessage] - start read loop\n")
	
	for {
		data, eob := mb.uartIn.ReadByte()
		
		// if we hit end of buffer before we find message return not found
		if eob != nil {
			log.Printf("[seekStartOfMessage] - return because we hit end of buffer")
			return false
		}
		
		// the '^' character is the start of a message
		if data == TOKEN_HAT {
			log.Printf("[seekStartOfMessage] - return because we hit start of message character")
			return true
		}
		
	}

}
