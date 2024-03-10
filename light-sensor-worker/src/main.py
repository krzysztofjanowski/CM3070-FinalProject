"""
A part of this light-sensor-worker script process has been taken from https://www.raspberrypi-spy.co.uk/?s=bh1750. 
This has been marked below. The part that comes from https://www.raspberrypi-spy.co.uk/?s=bh1750 is responsible for
communication with the light sensor and getting its readings. 

The rest has been written as part of CM3070 project and performs the following functions
1. Connect to MQTT broker  
2. Perform a light reading and publish it to the MQTT topic "light_sensor/light_reading"
3. The worker calculates absolute values between readings and publishes new reading only if there is significant difference between them, 
set to more than 1 light pixel but could be changed as needed 
"""

import smbus
import time
import paho.mqtt.client as mqtt

mqtt_address = "127.0.0.1"
# MQTT topic used to publish sensor readings 
topic = "light_sensor/light_reading"
frequency=1
previous_lightLevel = 0

# connect to mqtt 
client = mqtt.Client()
client.connect(mqtt_address)

# Code below comes from https://www.raspberrypi-spy.co.uk/?s=bh1750
#---------------------------------------------------------------------
#    ___  ___  _ ____
#   / _ \/ _ \(_) __/__  __ __
#  / , _/ ___/ /\ \/ _ \/ // /
# /_/|_/_/  /_/___/ .__/\_, /
#                /_/   /___/
#
#           bh1750.py
# Read data from a BH1750 digital light sensor.
#
# Author : Matt Hawkins
# Date   : 26/06/2018
#
# For more information please visit :
# https://www.raspberrypi-spy.co.uk/?s=bh1750
#
#---------------------------------------------------------------------

# Define some constants from the datasheet

DEVICE     = 0x23 # Default device I2C address

POWER_DOWN = 0x00 # No active state
POWER_ON   = 0x01 # Power on
RESET      = 0x07 # Reset data register value

# Start measurement at 4lx resolution. Time typically 16ms.
CONTINUOUS_LOW_RES_MODE = 0x13
# Start measurement at 1lx resolution. Time typically 120ms
CONTINUOUS_HIGH_RES_MODE_1 = 0x10
# Start measurement at 0.5lx resolution. Time typically 120ms
CONTINUOUS_HIGH_RES_MODE_2 = 0x11
# Start measurement at 1lx resolution. Time typically 120ms
# Device is automatically set to Power Down after measurement.
ONE_TIME_HIGH_RES_MODE_1 = 0x20
# Start measurement at 0.5lx resolution. Time typically 120ms
# Device is automatically set to Power Down after measurement.
ONE_TIME_HIGH_RES_MODE_2 = 0x21
# Start measurement at 1lx resolution. Time typically 120ms
# Device is automatically set to Power Down after measurement.
ONE_TIME_LOW_RES_MODE = 0x23

#bus = smbus.SMBus(0) # Rev 1 Pi uses 0
bus = smbus.SMBus(1)  # Rev 2 Pi uses 1

def convertToNumber(data):
  # Simple function to convert 2 bytes of data
  # into a decimal number. Optional parameter 'decimals'
  # will round to specified number of decimal places.
  result=(data[1] + (256 * data[0])) / 1.2
  return (result)

def readLight(addr=DEVICE):
  # Read data from I2C interface
  data = bus.read_i2c_block_data(addr,ONE_TIME_HIGH_RES_MODE_1)
  return convertToNumber(data)

# The end of the code from https://www.raspberrypi-spy.co.uk/?s=bh1750
#
#---------------------------------------------------------------------

# Code below is written by myself as part of CM3070 Final Project 

def main():
  global previous_lightLevel
  print("previous_lightLevel:", previous_lightLevel)
  # if there is no reading yet, we need to publish it so that dashboard can display it
  if previous_lightLevel == 0:
      lightLevel=readLight()
      message = "new light reading:" + str(format(lightLevel,'.2f'))
      client.publish(topic, message)
  while True:
    lightLevel=readLight() 
    # don't publish unless there is a significant difference in light conditions 
    if abs(lightLevel - previous_lightLevel)>1:
      previous_lightLevel = lightLevel
      message = "new light reading:" + str(format(lightLevel,'.2f'))
      client.publish(topic, message)
      print(message)
    time.sleep(frequency)
    

if __name__=="__main__":
   main()