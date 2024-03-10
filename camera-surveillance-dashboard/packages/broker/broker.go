package broker

import (
	"fmt"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/models"
	"log"
	"strconv"
	"strings"
    "math/rand"
    "time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT broker details
var brokerIp string = "192.168.178.64"
var portIp int = 1883
var Client mqtt.Client



var WebData models.WebData

// listens to MQTT messages and reacts accordingly.
var mqttMessagePublisherHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	messageReceived := string(msg.Payload())

    if strings.HasPrefix(messageReceived, "movement_detected-sensor1") {
        fmt.Println("Movement detected in sensor 1") 
        WebData.MovementSensor1 = true
    } else if strings.HasPrefix(messageReceived, "movement_not_detected-sensor1") {
        fmt.Println("No movement in sensor 1") 
        WebData.MovementSensor1 = false 
    }

    if strings.HasPrefix(messageReceived, "movement_detected-sensor2") {
        fmt.Println("Movement detected in sensor 2") 
        WebData.MovementSensor2 = true
    } else if strings.HasPrefix(messageReceived, "movement_not_detected-sensor2") {
        fmt.Println("No movement in sensor 2") 
        WebData.MovementSensor2 = false  
    }

    if strings.HasPrefix(messageReceived, "new light reading")  {
        //  to extract light intesity value, find index of first :
        colonIndex := strings.Index(messageReceived, ":")
        lightIntesityValueStr := messageReceived[colonIndex+1 : colonIndex+3]
        // convert to int 
        lightIntesityValueInt, err := strconv.Atoi(lightIntesityValueStr)
        if err != nil {
			log.Fatal("Error:", err)
			return
		}
        WebData.LightSensor = lightIntesityValueInt


    }
}

// to connect to MQTT broker  
var mqttConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    // MQTT client has connected successfully
}

// to handle disconnections 
var mqttConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("MQTT client connection has been lost with an error: %v", err)
}

func Broker() error {

    // MQTT client id 

    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())

    // Generate two random digits
    randomDigits := rand.Intn(100)

    // Append the random digits to the mqttClientId to make it somehow unique, needed for GET tests

    var mqttClientId string = fmt.Sprintf("camera-surveillance-backend%d", randomDigits)

    clientOptions := mqtt.NewClientOptions()
    clientOptions.AddBroker(fmt.Sprintf("tcp://%s:%d", brokerIp, portIp))
    clientOptions.SetClientID(mqttClientId)

    clientOptions.SetDefaultPublishHandler(mqttMessagePublisherHandler)
    clientOptions.OnConnect = mqttConnectHandler
    clientOptions.OnConnectionLost = mqttConnectLostHandler
    
	Client = mqtt.NewClient(clientOptions)
    if token := Client.Connect(); token.Wait() && token.Error() != nil {
        // panic(token.Error())
        return token.Error()
    }

	// Used to obtain sensor readings 
    error := subscribeToTopic(Client, "pir_sensor/motion_reading")
    if error != nil {
        return error
    }
	// Used to obtain sensor readings 
    error = subscribeToTopic(Client, "light_sensor/light_reading")
    if error != nil {
        return error
    }

    error = subscribeToTopic(Client, "camera/video")
    if error != nil {
        return error
    }

    return nil
    
}

func PublishToTopic(client mqtt.Client, topic string, message string) {

    token := client.Publish(topic, 0, false, message)
    token.Wait()
    time.Sleep(time.Second)
}

func subscribeToTopic(client mqtt.Client, mqttTopic string) error {

    topic := mqttTopic

    token := client.Subscribe(topic, 1, nil)
    if token := Client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    token.Wait()
  	// fmt.Printf("Subscribed to topic: %s", topic)
    return nil 
}
