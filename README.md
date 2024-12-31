# Zag-NetStats

**Zag-NetStats** is a lightweight and efficient network monitoring tool written in Go. It collects and displays real-time statistics about network interfaces, including data transfer speeds and total usage. With a clean interface and customizable output formats, it’s suitable for developers, system administrators, and anyone who needs detailed network statistics.

<br>
<div align="center">

[![Version](https://img.shields.io/github/v/release/ShadowZagrosDev/Zag-NetStats?label=Version&color=blue)](https://github.com/ShadowZagrosDev/Zag-NetStats/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/ShadowZagrosDev/Zag-NetStats/total?label=Downloads&color=success)](https://github.com/ShadowZagrosDev/Zag-NetStats/releases/latest)
[![Stars](https://img.shields.io/github/stars/ShadowZagrosDev/Zag-NetStats?style=flat&label=Stars&color=ff69b4)](https://github.com/ShadowZagrosDev/Zag-NetStats)

[![Go](https://img.shields.io/badge/Go-1.23.2-00ADD8.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ShadowZagrosDev/Zag-NetStats)](https://goreportcard.com/report/github.com/ShadowZagrosDev/Zag-NetStats)
[![Code Size](https://img.shields.io/github/languages/code-size/ShadowZagrosDev/Zag-NetStats?color=lightgrey)](https://github.com/ShadowZagrosDev/Zag-NetStats)
[![Top Language](https://img.shields.io/github/languages/top/ShadowZagrosDev/Zag-NetStats?color=yellowgreen)](https://github.com/ShadowZagrosDev/Zag-NetStats)

</div>
<br>


## Features

- **Real-Time Monitoring**: View instantaneous upload and download speeds.
- **Detailed Usage Statistics**: Monitor total data sent, received, and overall usage.
- **Customizable Output**:
  - JSON format for integration with other tools.
  - Tabular format for a clear and human-readable display.
- **Cross-Platform Support**: Works on Linux, macOS, and Windows.
- **Configurable Precision and Refresh Interval**: Fine-tune precision and update frequency as needed.

## Installation

### Option 1: Download from Releases

1. Go to the [Releases](https://github.com/ShadowZagrosDev/Zag-NetStats/releases) page of the repository.
2. Download the binary for your operating system.
3. Extract the downloaded file.
4. Run the binary:
   ```bash
   ./zag-netStats -i <interface_name>
   ```

### Option 2: Build from Source

1. Ensure you have Go installed (version 1.23.2 or higher).
2. Clone this repository:
   ```bash
   git clone https://github.com/ShadowZagrosDev/Zag-NetStats.git
   cd Zag-NetStats
   ```
3. Build the binary:
   ```bash
   go build -o zag-netStats ./cmd/main.go
   ```
4. Run the tool:
   ```bash
   ./zag-netStats -i <interface_name>
   ```


## Usage

### Command-Line Options

| Option          | Description                                       | Default Value |
| --------------- | ------------------------------------------------- | ------------- |
| `-i` (required) | Specify the network interface to monitor.         | N/A           |
| `-t`            | Refresh interval in seconds (1 to 3600).          | `1`           |
| `-p`            | Precision for rounding numerical values (0 to 6). | `2`           |
| `-f`            | Output format: `json` or `table`.                 | `table`       |

### Example

Monitor a network interface (`eth0`) with a refresh interval of 2 seconds and JSON output:

```bash
./zag-netStats -i eth0 -t 2 -f json
```


## Sample Output

### Tabular Format

```
+-----------+------------+------------+------------+------------+-------------+
| Interface | Sent Speed | Recv Speed | Total Sent | Total Recv | Total Usage |
+-----------+------------+------------+------------+------------+-------------+
| eth0      | 12.34 MB/s | 56.78 MB/s | 1.23 GB    | 4.56 GB    | 5.79 GB     |
+-----------+------------+------------+------------+------------+-------------+
```

### JSON Format

```json
{
  "interface": "eth0",
  "sentSpeed": { "value": 12.34, "unit": "MB/s" },
  "recvSpeed": { "value": 56.78, "unit": "MB/s" },
  "totalSent": { "value": 1.23, "unit": "GB" },
  "totalRecv": { "value": 4.56, "unit": "GB" },
  "totalUsage": { "value": 5.79, "unit": "GB" }
}
```


## How It Works

1. **Interface Selection**: The tool retrieves network I/O statistics for the specified interface using [gopsutil](https://github.com/shirou/gopsutil).
2. **Data Processing**:
   - Calculates instantaneous upload and download speeds.
   - Computes total data sent and received since the start of monitoring.
3. **Output Rendering**: Formats the data as a table or JSON for display.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


## Acknowledgments

- [gopsutil](https://github.com/shirou/gopsutil) for providing cross-platform system utilities.
- [tablewriter](https://github.com/olekukonko/tablewriter) for rendering tabular output.


---

**Made with ❤️ by ShadowZagrosDev**

