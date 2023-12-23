# mqtt2glados

This program allows you to play sound from [glados-tts](https://github.com/nerdaxic/glados-tts) using any linux (or bsd, but I didn't test that) computer and any music player that can play from stdin.

# running

First, you need to have glados-tts running. You can install it manually from [glados-tts](https://github.com/nerdaxic/glados-tts) or use [my pre-made docker image](https://github.com/LukaszMoskala/glados-tts-dockerfile).

Next, you need to compile this. When you have golang toolchain installed:
```
git clone https://github.com/LukaszMoskala/mqtt2glados.git
cd mqtt2glados.git
go build
```
And that's it. If you need to cross-compile this for diffrent architecture, you can just do `GOOS=linux GOARCH=arm64 go build`. For raspberry pi it will be either `arm64` or `arm`, depending on OS and hardware version.

After that, you need very simple `mqtt2glados.yaml`:
```yaml
mqtt_url: "tcp://homeassistant.local:1883"
mqtt_user: "automation"
mqtt_pass: "wow_much_secret"
mqtt_topic: "cmnd/mqtt2glados/test"
mqtt_client_id: "mqtt2glados-test"
glados_tts_endpoint: "http://127.0.0.1:8124/synthesize"
audio_player_command: "aplay"
audio_player_args: ["-"]
```
Keep in mind that glados_tts_endpoint needs to be reachable from machine running mqtt2glados. It can be the same machine or any other.

That's it, now any message published to `cmnd/mqtt2glados/test` will be synthesized and played using `aplay -`. You can test this using `mosquitto_pub -d -u automation -P wow_much_secret -h homeassistant.local -t "cmnd/mqtt2glados/test" -m "Hello there"`

You can use probably any MQTT broker, and any software to send messages (homeassistant, mosquitto_pub, node-red, whatever...). I didn't test any other music player, but aplay seems to work fine.

# Side notes

I put glados-tts docker image and this software in two days total. I assume that you run it on secure network or at least do IP-based firewalling.

My idea is to have GLaDOS say funny things when something happens in my room (like sensors detect motion when I'm away)

