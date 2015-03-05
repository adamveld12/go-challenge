#SPLICE spec

| offset  | Size(bytes) |Type           |
|:-------:|:---|:-----------------------|
| 0  | 13 | SPLICE header                |
| 13 | 1  | file size (bytes)             |
| 14 | 33 | hardware version size (bytes)|
| 46 | 4(float)| BPM                    |
| 50 | 4(int)  | first channel id       |
| 54 | 1(byte) | instrument name length |
| 55 | value@54*1(byte) | instrument name   |
| 54 + @55| 4(int) | pattern                |

Channel
4 bytes: id
1 byte: length of string
instrument: variable string
pattern: 16 bytes after instrument name
