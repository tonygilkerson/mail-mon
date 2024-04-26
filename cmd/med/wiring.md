# Wiring

| Pico Board Pin | Pico GPIO      | TFT Display w/buttons | Buzzer |
| -------------- | -------------- | --------------------- | ------ |
| 1              | GP0 (UART0 TX) |                       |        |
| 2              | GP1 (UART0 RX) |                       |        |
| 3              | GND            |                       |        |
| 4              | GP2            | dspKey2               |        |
| 5              | GP3            | dspKey3               |        |
| 6              | GP4            |                       |        |
| 7              | GP5            |                       |        |
| 8              | GND            |                       | GND    |
| 9              | GP6            |                       |        |
| 10             | GP7            |                       | Pos    |
| 11             | GP8            | dspDC                 |        |
| 12             | GP9            | dspCS                 |        |
| 13             | GND            |                       |        |
| 14             | GP10           | dspSCK                |        |
| 15             | GP11           | dspSDO                |        |
| 16             | GP12           | dspReset              |        |
| 17             | GP13           | dspLite               |        |
| 18             | GND            |                       |        |
| 19             | GP14           |                       |        |
| 20             | GP15           | dspKey0               |        |
| 21             | GP16           |                       |        |
| 22             | GP17           | dspKey1               |        |
| 23             | GND            |                       |        |
| 24             | GP18           |                       |        |
| 25             | GP19           |                       |        |
| 26             | GP20           |                       |        |
| 27             | GP21           |                       |        |
| 28             | GND            |                       |        |
| 29             | GP22           |                       |        |
| 30             | RUN            |                       |        |
| 31             | GP26           |                       |        |
| 32             | GP27           |                       |        |
| 33             | GND            |                       |        |
| 34             | GP28           | dspSDI                |        |
| 35             | ACD_VREF       |                       |        |
| 36             | 3v3 (out)      |                       |        |
| 37             | 3v3 (EN)       |                       |        |
| 38             | GND            |                       |        |
| 39             | 5v0 (VSYS)     |                       |        |
| 40             | 5v0 (VBUS)     |                       |        |

Not exposed as board pins

* **GP23** - OP Controls the on-board SMPS Power Save pin
* **GP24** - IP VBUS sense - high if VBUS is present, else low
* **GP25** - Onboard LED

## LED Buttons

```text
key2(top/left)----------------------key3(top/right)
|                                                 |
|                                                 |
|                                                 |
|                                                 |
|                                                 |
|                                                 |
|                                                 |
|                                                 |
key1(bottom/left)----------------key0(bottom/right)
```
 
* key0 - add 1hr to age
* key1 - add 30 min to age
* key2 - Reset/took meds
* key3 - subtract 1hr from age