package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
)

type ConfigT struct {
	MQTTURL            string   `yaml:"mqtt_url"`
	MQTTUser           string   `yaml:"mqtt_user"`
	MQTTPass           string   `yaml:"mqtt_pass"`
	MQTTTopic          string   `yaml:"mqtt_topic"`
	MQTTClientID       string   `yaml:"mqtt_client_id"`
	GladosTTSEndpoint  string   `yaml:"glados_tts_endpoint"`
	AudioPlayerCommand string   `yaml:"audio_player_command"`
	AudioPlayerArgs    []string `yaml:"audio_player_args"`
}

var conf ConfigT

func main() {
	f, err := os.Open("mqtt2glados.yaml")
	if err != nil {
		log.Fatalf("Failed to open config file: %s", err)
	}
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(conf.MQTTURL)
	opts.SetClientID(conf.MQTTClientID)
	opts.SetUsername(conf.MQTTUser)
	opts.SetPassword(conf.MQTTPass)
	msgchan := make(chan string, 32)
	opts.SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
		msgchan <- string(m.Payload())
	})
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %s", err)
	}
	token = client.Subscribe(conf.MQTTTopic, 2, nil)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("Failed to subscribe: %s", err)
	}
	log.Printf("Entering message loop")
	for {
		m := <-msgchan
		log.Printf("Received message: %s", m)

		URL := fmt.Sprintf("%s/%s", conf.GladosTTSEndpoint, url.PathEscape(m))

		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			log.Printf("Failed to create request: %s", err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to send request: %s", err)
			continue
		}
		defer resp.Body.Close()
		proc := exec.Command(conf.AudioPlayerCommand, conf.AudioPlayerArgs...)
		proc.Stdin = resp.Body
		err = proc.Run()
		if err != nil {
			log.Printf("Failed to run audio player: %s", err)
		}

	}

}
