# Wiring

| Pico Board Pin | Pico GPIO      | e-Paper Dsp | Buttons                           | Com Pico       |neo
| -------------- | -------------- | ----------- | --------------------------------- | -------------- |----
| 1              | GP0 (UART0 TX) |             |                                   | GP5 (UART1 RX) |
| 2              | GP1 (UaRT0 RX) |             |                                   | GP4 (UART1 TX) |
| 3              | GND            |             |                                   |                |
| 4              | GP2            |             | MbxDoorOpenedAck PinInputPulldown |                |
| 5              | GP3            |             | Request          PinInputPulldown |                |
| 6              | GP4            |             |                                   | GP1            |
| 7              | GP5            |             |                                   | GP0            |
| 8              | GND            |             |                                   |                |
| 9              | GP6            |             |                                   |                | DIN
| 10             | GP7            |             |                                   |                |
| 11             | GP8            |             |                                   |                |
| 12             | GP9            |             |                                   |                |
| 13             | GND            |             |                                   |                |
| 14             | GP10           |             |                                   |                |
| 15             | GP11           | DC          |                                   |                |
| 16             | GP12           | RST         |                                   |                |
| 17             | GP13           | BUSY        |                                   |                |
| 18             | GND            |             |                                   |                |
| 19             | GP14           |             |                                   |                |
| 20             | GP15           |             |                                   |                |
| 21             | GP16           |             |                                   |                |
| 22             | GP17           | CS          |                                   |                |
| 23             | GND            |             |                                   |                |
| 24             | GP18           | CLK         |                                   |                |
| 25             | GP19           | DIN         |                                   |                |
| 26             | GP20           |             |                                   |                |
| 27             | GP21           |             |                                   |                |
| 28             | GND            |             |                                   |                |
| 29             | GP22           |             |                                   |                |
| 30             | RUN            |             |                                   |                |
| 31             | GP26           |             |                                   |                |
| 32             | GP27           |             |                                   |                |
| 33             | GND            |             |                                   |                |
| 34             | GP28           |             |                                   |                |
| 35             | ACD_VREF       |             |                                   |                |
| 36             | 3v3 (out)      | VIN         |                                   |                |
| 37             | 3v3 (EN)       |             |                                   |                |
| 38             | GND            | GND         |                                   |                | GND
| 39             | 5v0 (VSYS)     |             |                                   |                | 5v
| 40             | 5v0 (VBUS)     |             |                                   |                |

Not exposed as board pins

* **GP23** - OP Controls the on-board SMPS Power Save pin
* **GP24** - IP VBUS sense - high if VBUS is present, else low
* **GP25** - Onboard LED
