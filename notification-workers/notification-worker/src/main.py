
import sys
import time 
import paho.mqtt.client as mqtt
sys.path.append('../../') # to allow import of below packages 
from smsworker.src.main import SmsWorker
from slackworker.src.main import SlackWorker
from emailworker.src.main import EmailWorker
import datetime
import json


# registraionDetailsData file path
file_path = "/opt/camera-surveillance-system/camera-surveillance-dashboard/registraionDetailsData.json"

# MQTT Broker runs on the pi hence local address 
mqtt_address = "127.0.0.1"
# MQTT topic that the notification worker subscribes to and listens for new videos messages arrival 
topic = "camera/video"
password = ""
notification_message = "Your Camera Surveillance system has recorded a new video. It has been published to the web dashboard and is ready to view and download."

# global config variables 
phone_number = ""
api_key = ""
destination_email = ""

with open(file_path, 'r') as registraion_details:
    data = json.load(registraion_details)
    phone_number = data["PhoneNumber"]
    api_key = data["SlackKey"]
    destination_email = data["Email"]
 
sms_to_number = phone_number
# this number is fixed, relies on API gateway, here registration with Twilio
sms_from_number = "+18587035225"

# configure sms object  
sms_worker = SmsWorker()

# configure email object and send email 
email_worker = EmailWorker()

# configure slack object and send slack message
slack_worker = SlackWorker()

def on_message(client, userdata, message):
    message_received = str(message.payload.decode("utf-8"))
    print("Message received:",  message_received)

    if message_received == "new_processed_video":
        print("A new processed video arrived, I am going to send notifiation.")  
  
        # here comes the action of sending emails, sms etc 
        print("going to send an email notification...")
        email_worker.send_email(contents=notification_message, dst_email=destination_email)

        time.sleep(0.5)
        print("going to send an sms notification...")
        sms_worker.send_sms(message_content=notification_message, from_number=sms_from_number, to_number=sms_to_number)

        time.sleep(0.5)
        print("going to send an slack notification...")
        slack_worker.send_status_to_slack(api_token=api_key, message=notification_message)

        print("going to write to notification log...")
        with open("/opt/camera-surveillance-system/camera-surveillance-dashboard/notifications.txt", "a") as file:
            # Write to the notification log file
            current_datetime = datetime.datetime.now()
            formatted_datetime = current_datetime.strftime("%Y-%m-%d|%H:%M:%S")
            message = "A new notification sent to sms, slack and email.\n"
            message_content = f"{formatted_datetime}: {message}"
            file.write(message_content)
            

def main():
    print("executing main notification worker")

    time.sleep(1)
    client = mqtt.Client()
    client.on_message = on_message

    client.connect(mqtt_address)

    client.subscribe(topic)

    client.loop_forever()

    
if __name__ == "__main__":
    main()