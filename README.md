# GeoIP Seeker Server

> The geoip seeker server can be used to query ip geography

# Supported

- [x] CZ88.NET QQWay
- [x] IPIP.NET DAT
- [x] IPIP.NET DATX
- [ ] MaxMind GeoIP2
- [ ] MaxMind GeoLite2

# Deployment

1. usage `build.sh` build binary files
2. copy `build` directory to `/opt/geoip-seeker-server`
3. copy `build/geoip-seeker-server.service` file to `/etc/systemd/system`
4. enable service `systemd enable geoip-seeker-server`
5. start service `systemd start geoip-seeker-server`
