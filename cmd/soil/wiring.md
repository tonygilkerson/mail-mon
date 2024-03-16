# Wiring

| Pico Board Pin | Pico GPIO      | Lora Breakout Board | Charger | Soil Sensor | 7-segment display | L293D          | Solenoid   | 9v Battery   |
| -------------- | -------------- | ------------------- | ------- | ----------- | ----------------- | -------------- | ---------- | ------------ |
| 1              | GP0 (UART0 TX) |                     |         |             |                   |                |            |              |
| 2              | GP1 (UART0 RX) |                     |         |             |                   |                |            |              |
| 3              | GND            |                     |         |             |                   | Pin-4 (GND)    |            |              |
| 4              | GP2            |                     |         |             |                   |                |            |              |
| 5              | GP3            |                     |         |             |                   |                |            |              |
| 6              | GP4            |                     |         |             |                   |                |            |              |
| 7              | GP5            |                     |         |             |                   |                |            |              |
| 8              | GND            |                     |         |             |                   | Pin-5 (GND)    |            |              |
| 9              | GP6            |                     |         |             |                   | Pin-1 (Enable) |            |              |
| 10             | GP7            |                     |         |             |                   | Pin-2 (In1)    |            |              |
| 11             | GP8            |                     |         |             |                   | Pin-7 (In2)    |            |              |
| 12             | GP9            |                     |         |             |                   |                |            |              |
| 13             | GND            |                     |         |             | GND (black)       |                |            |              |
| 14             | GP10           |                     |         |             | CLK (yellow)      |                |            |              |
| 15             | GP11           |                     |         |             | DIO (white)       |                |            |              |
| 16             | GP12           |                     |         | SDA(White)  |                   |                |            |              |
| 17             | GP13           |                     |         | SCL(Green)  |                   |                |            |              |
| 18             | GND            |                     |         | GND (Black) |                   |                |            |              |
| 19             | GP14           |                     |         |             |                   |                |            |              |
| 20             | GP15           | EN                  |         |             |                   |                |            |              |
| 21             | GP16           | MISO (SPI0)         |         |             |                   |                |            |              |
| 22             | GP17           | CS                  |         |             |                   |                |            |              |
| 23             | GND            |                     | GND     |             |                   |                |            |              |
| 24             | GP18           | SCK  (SPI0)         |         |             |                   |                |            |              |
| 25             | GP19           | MOSI (SPI0)         |         |             |                   |                |            |              |
| 26             | GP20           | RST                 |         |             |                   |                |            |              |
| 27             | GP21           | G0                  |         |             |                   |                |            |              |
| 28             | GND            |                     |         |             |                   |                |            |              |
| 29             | GP22           | G1                  |         |             |                   |                |            |              |
| 30             | RUN            |                     |         |             |                   |                |            |              |
| 31             | GP26           |                     |         |             |                   |                |            |              |
| 32             | GP27           |                     |         |             |                   |                |            |              |
| 33             | GND            |                     |         |             |                   |                |            |              |
| 34             | GP28           |                     |         |             |                   |                |            |              |
| 35             | ACD_VREF       |                     |         |             |                   |                |            |              |
| 36             | 3v3 (out)      | VIN                 |         | VIN (Red)   | VIN (red)         |                |            |              |
| 37             | 3v3 (EN)       |                     |         |             |                   |                |            |              |
| 38             | GND            | GND                 |         |             |                   | Pin-13 (GND)   |            | Neg terminal |
| 39             | 5v0 (VSYS)     |                     | 5V      |             |                   | Pin-16 (VSS)   |            |              |
| 40             | 5v0 (VBUS)     |                     |         |             |                   |                |            |              |
|                |                |                     |         |             |                   | Pin-8 (VSmot)  |            | Pos terminal |
|                |                |                     |         |             |                   | Pin-3 (Out1)   | Terminal A |              |
|                |                |                     |         |             |                   | Pin-6 (Out2)   | Terminal B |              |

Not exposed as board pins

* **GP23** - OP Controls the on-board SMPS Power Save pin
* **GP24** - IP VBUS sense - high if VBUS is present, else low
* **GP25** - Onboard LED
