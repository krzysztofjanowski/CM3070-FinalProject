"""
Camera worker key functions are:
1. Connect to MQTT broker and subscribe to "pir_sensor/motion_reading" topic
2. When the "pir_sensor/motion_reading" topic receives a message "movement_detected", record a video of length set in the 'duration' parameter
3. Save the video as "/opt/unprocessed_videos/raw.mp4"
4. Publish a "new_video" message in the "camera/video" topic to triger further video processing from other services
It also contains a special privacy flag which when set to true will prevent recording of any videos.
"""
import paho.mqtt.client as mqtt
import picamera2
import os 
import time


# mqtt topic used for notification about movement detection 
topic_motion = "pir_sensor/motion_reading"

# MQTT Broker address which could be local or remote depending on architecture, here left to loopback as an example 
mqtt_broker_address = "127.0.0.1"

# Destination for the raw videos 
video_file_destination = "/opt/unprocessed_videos/raw.mp4"

# Create the destination directory if it doesn't exist
os.makedirs("/opt/unprocessed_videos/", exist_ok=True)

# main picamera2 object 
camera = picamera2.Picamera2()

# privacy_flag to disable recording
privacy_flag = False 

# mqtt topic that camera-worker subscribes to for notification about privacy 
topic_video = "camera/video"

# flat to prevent multiple video recording requests
video_recording_in_progress = False 

def on_message(client, userdata, message):
    """
    determines what to do when a particular message of interest arrives
    :param client: mqtt client 
    :param userdata: 
    :param message: message that is received 
    :return: True if all went successful, False if there was an issue 
    """
    global privacy_flag, video_recording_in_progress
    message_received = str(message.payload.decode("utf-8"))
    print("Message received:",  message_received)

    # sensor 1 and sensor 2 publish movement_detected-sensor1 and movement_detected-sensor2 respectively hence startswith
    if message_received.startswith("movement_detected") and privacy_flag == False and video_recording_in_progress == False:
        video_recording_in_progress = True 
        recording_result = record_video(horizontal_flip=True)
        if recording_result == False:
            return False 
        print("Recording has finished.")
        client.publish("camera/video", "new_video")
        video_recording_in_progress = False 
        # sleep so that the image processing service can finish and is not overwhelmend or raw.mp4 is overwritten too early 
        time.sleep(10)
        return True 
    elif message_received.startswith("movement_detected") and  privacy_flag == True:
        print("Movement detected but privacy flag set to true so no new video is going to recorded")
    elif message_received.startswith("movement_detected") and  video_recording_in_progress == True:
        print("Movement detected but video recording in progress so no new video is going to recorded")  

    if message_received == "DISABLE ALL RECORDING":
        privacy_flag = True
        print("privacy_flag changed to: ", privacy_flag)
    elif message_received == "ENABLE ALL RECORDING":
        privacy_flag = False 
        print("privacy_flag changed to: ", privacy_flag)


def record_video(horizontal_flip):
    """
    records the video
    :param horizontal_flip: if camera is upside down, video needs to be flipped
    :return: True if all went successful, False if there was an issue 
    """
    global camera

    camera.resolution = (640, 480)

    # Enable horizontal flip
    camera.vflip = horizontal_flip
    try:
        print("starting recording of a new video...")
        camera.start_and_record_video(video_file_destination, duration=5)
    except:
        print("record_video:: an exception occurred")
        return False 
    return True 


def main():
    # MQTT client, address, topic and message
    client = mqtt.Client()

    MQTT_TOPIC = [(topic_motion,0),(topic_video,0)]

    client.connect(mqtt_broker_address)
    client.subscribe(MQTT_TOPIC)

    # this part will trigger video capture on 'movement detected' message
    client.on_message = on_message
    print("privacy_flag:", privacy_flag)
    # loop with slight delay to avoid overutilziation of cpu
    time.sleep(1)
    client.loop_forever()


if __name__ == "__main__":
    main()