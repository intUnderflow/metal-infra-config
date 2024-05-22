# Metal Infra Config

Goals:
- Shared configuration across peers
- Authentication
-Hosted on each peer, totally resilient to downtime

Used to:
- Support bootstrapping a valid Wireguard mesh network configuration
- Peers on the mesh can have dynamic IPs

Design:
- Only one peer can write to one key
- All values are strings
- Peers increase vector clocks when they write keys
- Peers sync keys to each other
