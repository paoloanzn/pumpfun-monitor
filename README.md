# Pump.fun Monitor

Real-time monitoring system for Pump.fun token mints and migrations on Solana blockchain. Built for high scalability using producer-consumer architecture with WebSocket streaming.

[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-green.svg)](https://opensource.org/licenses/)

## Key Features ✨
- **Dual Monitoring** - Track both new token mints AND liquidity migrations
- **Scalable Architecture** - Multiple worker processes with message queues

## System Architecture 🏗
[![](https://mermaid.ink/img/pako:eNp9ks1ugzAQhF_F8pkcypFDpQbyWxGpIlKl4hxcvAWLsEbGTlUlefcaKCStmtxmR59nV9490kwJoAHNNa8Lso0YPqVbVQKSWKIhsUJplN6RyeTxtIFPUrUuHABNcyLTtIegaXgO5MWChR3D0Nkuz0iFvwMu9pAQXaN_YqbtGzLre7zCe6KyEgxJQB9Ak4eRmN8g_JFY3CA2jog6Ynk1x7-temx1D_NHbH0Pa5vOOuw5Dfey_QfnzHuH4WIQy0GsBrEeBEPq0Qp0xaVwuzsyJIRRU0AFjAZOCq5LRhmeHcetUckXZjQw2oJHtbJ5QYMPvm9cZWvBDUSSu3mr0a05vil1qUG0S4z7U-kuxnMX0zb_yQQUoENl0dDAP38DcUPH9w?type=png)](https://mermaid.live/edit#pako:eNp9ks1ugzAQhF_F8pkcypFDpQbyWxGpIlKl4hxcvAWLsEbGTlUlefcaKCStmtxmR59nV9490kwJoAHNNa8Lso0YPqVbVQKSWKIhsUJplN6RyeTxtIFPUrUuHABNcyLTtIegaXgO5MWChR3D0Nkuz0iFvwMu9pAQXaN_YqbtGzLre7zCe6KyEgxJQB9Ak4eRmN8g_JFY3CA2jog6Ynk1x7-temx1D_NHbH0Pa5vOOuw5Dfey_QfnzHuH4WIQy0GsBrEeBEPq0Qp0xaVwuzsyJIRRU0AFjAZOCq5LRhmeHcetUckXZjQw2oJHtbJ5QYMPvm9cZWvBDUSSu3mr0a05vil1qUG0S4z7U-kuxnMX0zb_yQQUoENl0dDAP38DcUPH9w)
## Installation 📦

### From Source
**Prerequisites**:
- Go 1.24.1+
- Node.js 18+ (for benchmarking)

```bash
git clone https://github.com/paoloanzn/pumpfun-monitor.git
cd pumpfun-monitor
make install
```

## Usage 🚦

### Start Monitoring Service
```bash
pumpfun-monitor start -mint-workers 3 -migration-workers 2 -max-recon 100
```

**Flags**:
- `-mint-workers`: Number of mint monitoring processes
- `-migration-workers`: Migration monitoring processes
- `-max-recon`: Maximum reconnection attempts

### WebSocket Endpoints
- **Mint Stream**: `ws://localhost:8080/{mint-uuid}`
- **Migration Stream**: `ws://localhost:8080/{migration-uuid}`

## API Reference 📡

### HTTP Endpoints
**Get Active Consumers**:
Get mint consumers
```bash
curl http://localhost:8080/mint-consumers
```

Get migration consumers
```bash
curl http://localhost:8080/migration-consumers
```

**Response**:
```output
{
"uuids": ["550e8400-e29b-41d4-a716-446655440000"]
}
```

## License 📄
This project is licensed under the GNU GPLv3 License - see [LICENSE](LICENSE) file for details.

## Contributing 🤝
PRs welcome! Please follow existing code patterns and include tests where applicable.