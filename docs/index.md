# MBX-IOT

This repo contains the software used for my mailbox IoT project, it will:

* Detect when the U.S. mail has been delivered (when the mailbox door is opened)
* Detect when a car passes by the mailbox
* Once detected a message is sent via a radio signal to my house where it is recorded and notifications sent

![mbx-gateway](img/mbx-gateway.drawio.png)

---
There are two apps (TinyGo binaries) contained in this repo

**mbx** - This app is loaded on a Raspberry Pi Pico along with several peripheral devices. The unit is contained in a weather proof electrical junction box mounted to my mailbox. It perform the following functions:

* detect when the mailbox door is opened
* detect when a car passes by
* charge the battery from a solar cell

**gateway** - This app is loaded on a LoRa-E5 Dev Board that is connected to a server located in my the house. The server hosts a K8s cluster running monitoring and alerting software. The **gateway** performs the following functions:

* receive radio signals from the **mbx** device
* sends the message to an application running in the K8s cluster via serial port
* once the messages are received in the K8s cluster they are stored in a time series db so that graphical dashboards and be produced and notifications sent

## Devices

### Pico

![pico pins](img/pico-pins.png)

* [pico](https://www.adafruit.com/product/4864)

### LoRa Breakout Board

![lora](img/lora-breakout-board.png)

* [Adafruit RFM95W LoRa Radio Transceiver Breakout - 868 or 915 MHz - RadioFruit](https://www.adafruit.com/product/3072)

### Solar Charger

![charger](img/Adafruit-charger-bq24074.jpg)

* [Adafruit Universal USB / DC / Solar Lithium Ion/Polymer charger - bq24074](https://www.adafruit.com/product/4755)

![charger](img/solar-panel.png)

* [Small 6V 1W Solar Panel](https://www.adafruit.com/product/3809)

### LoRa E5 Dev Board

![lora-e5-pins](img/lora-e5-dev-kit-pins.jpg)

* [Wio-E5 Dev Kit - STM32WLE5JC, ARM Cortex-M4 and SX126x embedded, supports LoRaWAN on EU868 & US915](https://www.seeedstudio.com/LoRa-E5-Dev-Kit-p-4868.html)

### Battery

![batt](img/batt.jpeg)

* [Lithium Ion Cylindrical Battery - 3.7v 2200mAh](https://www.adafruit.com/product/1781)

### Photo cell

![photo cell](img/photo-cell.jpeg)
* [Photo cell (CdS photoresistor)](https://www.adafruit.com/product/161)

### Transistors

![Transistors](img/transistor.jpeg)

* [NPN Bipolar Transistors (PN2222) - 10 pack
](https://www.adafruit.com/product/756)
