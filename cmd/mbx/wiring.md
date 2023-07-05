# Wiring

DEVTODO - This needs verified

| Pico     | Lora Breakout Board | Charger Breakout Board | Solar Cell | Battery | Photo Cell w/Transistor on Mailbox Door | Photo Cell w/Transistor on Mule Tx LED |
| -------- | ------------------- | ---------------------- | ---------- | ------- | --------------------------------------- | -------------------------------------- |
| 3v3(out) | VIN                 | VBUS                   | Red        |         | T2                                      | T2                                     |
| GND      | GND                 | GND                    | Black      | Black   |                                         |                                        |
| GP10     |                     | CHG                    |            |         |                                         |                                        |
| GP11     |                     | PGOOD                  |            |         |                                         |                                        |
| GP12     |                     |                        |            |         |                                         | T1                                     |
| GP13     |                     |                        |            |         | T1                                      |                                        |
| GP15     | EN                  |                        |            |         |                                         |                                        |
| GP16     | MISO                |                        |            |         |                                         |                                        |
| GP17     | CS                  |                        |            |         |                                         |                                        |
| GP18     | SCK                 |                        |            |         |                                         |                                        |
| GP19     | MOSI                |                        |            |         |                                         |                                        |
| GP20     | RST                 |                        |            |         |                                         |                                        |
| GP21     | G0                  |                        |            |         |                                         |                                        |
| GP22     | G1                  |                        |            |         |                                         |                                        |
|          |                     | LIPO                   |            | Red     |                                         |                                        |
