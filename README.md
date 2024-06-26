# pc2mqtt

Control your PC or homeserver with homeassistant or any other MQTT enabled home automation system.

## Usage

pc2mqtt run on your pc or homeserver and exposes its state and actions via MQTT.

## Getting Started / Installation

> Installation of pc2mqtt works a little different on windows vs linux.

### Linux

1. Download the latest binary from the releases
2. Setup a systemd service running the binary on start
3. Thats it

### Windows

For windows pc2mqtt uses the [windows-service-wrapper](https://github.com/winsw/winsw).

1. Download the latest windows zip archive from the releases
2. Unzip it to dome directory of your choice
3. In cmd run `pc2mqtt.exe install` and `pc2mqtt.exe start` to install and start it as a windows service

## Config

When first starting the application, a `config.json` will be created right next to it. It looks like this:
´´´json
{
    "device_id": "63fbeebb-f107-4903-ab36-6104b9d802b0",
    "device_name": "MY-PC-HOSTNAME",
    "mqtt": {
        "host": "<YOUR MQTT HOST>",
        "port": 1883,
        "username": "<\u003cMQTT USER>",
        "password": "<\u003cMQTT PASSWORD>",
        "auto_discovery_prefix": "homeassistant"
    },
    "debug_mode": false
}
´´´

- device_id: a generated id to identify your device. Can be changed to any unique value of your choice
- deviec_name: How your device will be named in eg. homeassistant. Defaults to your hostname, can be changed to any value
- mqtt.host: Your mqtt hostname eg. 192.168.0.10
- mqtt.port: Your mqtt port
- mqtt.username: Your mqtt username
- mqtt.password: Your mqtt password
- mqtt.auto_discovery_prefix: The prefix used for the auto discovery messages. Defaults to `homeassistant`. Can be changed to fit your needs
