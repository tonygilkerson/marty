# MBX-IOT

This repo contains the software used for my mailbox IoT project, it will:

* Detect when the U.S. mail has been delivered (when the mailbox door is opened)
* Detect when a car passes by the mailbox
* Once detected a message is sent via a radio signal to my house where it is recorded and notifications sent

![mbx-gateway](img/mbx-gateway.drawio.png)

---
There are two apps (TinyGo binaries) contained in this repo

* **mbx** - This app is loaded on a Raspberry Pi Pico along with several peripheral devices. The unit is contained in a weather proof electrical junction box mounted to my mailbox. It perform the following functions:
  * detect when the mailbox door is opened
  * detect when a car passes by
  * charge the battery from a solar cell
* **gateway** - This app is loaded on a LoRa-E5 Dev Board that is tethered a server located my the house. The server hosts a K8s cluster running monitoring and alerting software. The **gateway** performs the following functions:
  * receive radio signals from the **mbx** device
  * sends the message to an application running in the K8s cluster via serial port
  * once the messages are received in the K8s cluster they are stored in a time series db so that graphical dashboards and be produced and notifications sent

## Devices

![pico pins](img/pico-pins.png)