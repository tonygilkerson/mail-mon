# Wiring

| Pico      | Lora Breakout Board | e-Paper Dsp    | NeoPixel Stick X 8 |
| --------- | ------------------- | -------------- | ------------------ |
| 5v0(VSYS) |                     |                | 5VDC               |
| 3v3(out)  | VIN                 | VIN            |                    |
| GND       | GND                 | GND            |                    |
| GP0       |                     |                |                    |
| GP1       |                     |                |                    |
| GP2       |                     |                | DIN                |
| GP3       |                     |                |                    |
| GP4       |                     |                |                    |
| GP5       |                     |                |                    |
| GP6       |                     |                |                    |
| GP7       |                     | *[DC]          |                    |
| GP8       |                     | *[RX SPI1]     |                    |
| GP9       |                     | *[CS SPI1]     |                    |
| GP10      |                     | CS *[SCK SPI1] |                    |
| GP11      |                     | DC *[TX SPI1]  |                    |
| GP12      |                     | RST            |                    |
| GP13      |                     | BUSY           |                    |
| GP14      |                     |                |                    |
| GP15      | EN                  |                |                    |
| GP16      | MISO (SPI0)         |                |                    |
| GP17      | CS                  |                |                    |
| GP18      | SCK  (SPI0)         | CLK (SPI0)     |                    |
| GP19      | MOSI (SPI0)         | DIN (SPI0)     |                    |
| GP20      | RST                 |                |                    |
| GP21      | G0                  |                |                    |
| GP22      | G1                  |                |                    |
| GP23      |                     |                |                    |
| GP26      |                     |                |                    |
| GP27      |                     |                |                    |
| GP28      |                     |                |                    |
