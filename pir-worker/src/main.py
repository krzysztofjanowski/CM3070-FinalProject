"""
PIR (aka movement) worker script key functions are:
1. Connect to MQTT broker 
2. Check readings from the movement sensor (if the reading is set to 'High' then movement occur)
3. Publish a "movement_detected-sensor1" message in the "pir_sensor/motion_reading" topic 
"""

import RPi.GPIO as GPIO
import time
import paho.mqtt.client as mqtt


GPIO.setmode(GPIO.BCM)

# pin that PIR sensor is connected to 
PIR_PIN=27
# the length of pause if movement occurs as there is no point in constant readings 
sleep_interval = 10

# make GPIO PIR as input, get 0/low (default) or 1/high if movement occurs
GPIO.setup(PIR_PIN,GPIO.IN, pull_up_down=GPIO.PUD_DOWN)

# MQTT Broker address
mqtt_address = "127.0.0.1"
# MQTT topic that the pir worker uses to send pir sensor message if movement occured 
topic = "pir_sensor/motion_reading"
message = "movement_detected-sensor1"

# connect to MQTT broker
client = mqtt.Client()
client.connect(mqtt_address)

def main():
    old_pir_reading = pir_reading = GPIO.input(PIR_PIN)
    while True:
        # to slow down the loop, otherwise cpu utilization goes very high 
        time.sleep(1)
        pir_reading = GPIO.input(PIR_PIN)
        if pir_reading == "0" or pir_reading == 0:
            print("Movement detected: false")
        else:
            print("Movement detected: true")            
        if pir_reading == 1 and pir_reading != old_pir_reading:
            print("movement detected sensor1")
            client.connect(mqtt_address)
            client.publish(topic, message)
            old_pir_reading = pir_reading
            time.sleep(sleep_interval)
        elif pir_reading == 0 and pir_reading != old_pir_reading:
            print("no movement sensor2") 
            client.connect(mqtt_address)
            client.publish(topic, "movement_not_detected-sensor1")
            time.sleep(1)
            old_pir_reading = pir_reading


if __name__ == "__main__":
    main()