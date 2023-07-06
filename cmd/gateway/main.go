package main

import (
	"log"
	"machine"
	"runtime"
	"strings"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
	"tinygo.org/x/drivers/sx127x"
)

const (
	// HEARTBEAT_DURATION_SECONDS = 300
	// TXRX_LOOP_TICKER_DURATION_SECONDS = 300
	HEARTBEAT_DURATION_SECONDS = 60
	TXRX_LOOP_TICKER_DURATION_SECONDS = 30
)


/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	//
	// Named PINs
	//
	var en machine.Pin = machine.GP15
	var sdi machine.Pin = machine.GP16 // machine.SPI0_SDI_PIN
	var cs machine.Pin = machine.GP17
	var sck machine.Pin = machine.GP18 // machine.SPI0_SCK_PIN
	var sdo machine.Pin = machine.GP19 // machine.SPI0_SDO_PIN
	var rst machine.Pin = machine.GP20
	var dio0 machine.Pin = machine.GP21 // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	var dio1 machine.Pin = machine.GP22 // (GP22--G1)I don't now what this does but it seems to need to be connected
	var uartTx machine.Pin = machine.GP0 // machine.UART0_TX_PIN
	var uartRx machine.Pin = machine.GP1 // machine.UART0_RX_PIN
	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 10)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// setup Uart
	//
	machine.UART0.Configure(machine.UARTConfig{BaudRate: 115200, TX: uartTx, RX: uartRx})


	//
	// 	Setup Lora
	//
	var loraRadio *sx127x.Device
	// I am thinking that a batch of message can be half dozen max so 250 should be plenty large
	txQ := make(chan string, 250) 
	rxQ := make(chan string, 250) 

	radio := road.SetupLora(*machine.SPI0, en, rst, cs, dio0, dio1, sck, sdo, sdi, loraRadio, &txQ, &rxQ, 0, 0, TXRX_LOOP_TICKER_DURATION_SECONDS, road.TxRx)

	// Launch go routines
	go writeToSerial(&rxQ, machine.UART0)
	go radio.LoraRxTx()
	

	// Main loop
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int
	
	for range ticker.C {
		
		log.Printf("------------------MainLoopHeartbeat-------------------- %v", count)
		count += 1

		//
		// Send Heartbeat to Tx queue
		//
		txQ <- "GatewayMainLoopHeartbeat"

		runtime.Gosched()
	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func runLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
		print("run-")
	}

}

func writeToSerial(rxQ *chan string, uart *machine.UART) {
	var msgBatch string

	for msgBatch = range *rxQ {
			
		log.Printf("Message batch: [%v]", msgBatch)
		
		messages := strings.Split(string(msgBatch), "|")
		for _, msg := range messages {
			log.Printf("Write to serial: [%v]", msg)
			uart.Write([]byte(msg))
			time.Sleep(time.Millisecond * 50) // Mark the End of a message
		}
	
		runtime.Gosched()

	}

}



