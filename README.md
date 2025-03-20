# Pumpfun Monitor

Tool to monitor pumpfun token launches and token migrations.
All in real-time without any API key required.

Its built to be easily used at scale, designed to be extended with ease using a simple *data-producer*/*data-consumer* paradigm.

[![](https://mermaid.ink/img/pako:eNp1kEFrwzAMhf-K0WmD9hD3lkMv3XGGtrstzkHEamLaWMWxGaPtf5_SdAsMqoPQ86dneLpAw46ghMOJv5oOY1LvexuU1LaojA_Jh1YZDj5xrNVyub4aGgZs6arMrqheHkrtMmV6rR9WLdY2YvIcnpr1E_PU5fPRoDZFteEw5J6iKupfpiemZ6b_s9XMVsJgATL26J2EvYy7FlJHPVkoZXQYjxZsuMke5sQf36GBMsVMC4ic2w7KA54GUfnsMNGbR8nX_72eMXwyz5rcmNlMt72f-PYDdjp0ww?type=png)](https://mermaid.live/edit#pako:eNp1kEFrwzAMhf-K0WmD9hD3lkMv3XGGtrstzkHEamLaWMWxGaPtf5_SdAsMqoPQ86dneLpAw46ghMOJv5oOY1LvexuU1LaojA_Jh1YZDj5xrNVyub4aGgZs6arMrqheHkrtMmV6rR9WLdY2YvIcnpr1E_PU5fPRoDZFteEw5J6iKupfpiemZ6b_s9XMVsJgATL26J2EvYy7FlJHPVkoZXQYjxZsuMke5sQf36GBMsVMC4ic2w7KA54GUfnsMNGbR8nX_72eMXwyz5rcmNlMt72f-PYDdjp0ww)

![](/assets/Screenshot%202025-03-20%20at%2001.01.25.png)

## Installation

_Pre compiled binaries to be added soon._

## Manual Installation
### Prerequisites
- Go >= go1.24.1
- Git
- Make
```bash
git clone https://github.com/paoloanzn/pumpfun-monitor.git && cd pumpfun-monitor && make install
```

## Usage

To start the monitors and the logging run:
```bash
pumpfun-monitor start
```