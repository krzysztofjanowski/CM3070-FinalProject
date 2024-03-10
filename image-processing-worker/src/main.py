"""
Image processing worker key functions are:
1. Connect to MQTT broker and subscribe to "camera/video" topic
2. When the "camera/video" topic receives a message "new_processed_video", transfer the raw.mp4 video from edge to fog layer
3. Process the video using openCV to detect, track and mark movement
4. Transform the video using ffmpeg so that it is ready for viewing in the browser
5. Publish a "new_processed_video" message in the "camera/video" topic 
"""

import time
import paho.mqtt.client as mqtt
import cv2
import os 
from datetime import datetime
import subprocess
import subprocess


# MQTT client, address, topic and message
mqtt_address = "127.0.0.1"
topic = "camera/video"
message = "new_processed_video"

video_source_file = '/opt/unprocessed_videos/raw.mp4'
destination_dir = '/opt/processed_videos/'

# Create the destination directory if it doesn't exist
os.makedirs(destination_dir, exist_ok=True)

def on_message(client, userdata, message):
    message_received = str(message.payload.decode("utf-8"))
    print("Message received:",  message_received)

    if message_received == "new_video":
        
        # generate current date and time stirng in form of 20/02/2024-05:23:30 etc
        now = datetime.now()
        date_time = now.strftime("%Y-%m-%d %H:%M:%S")

        print(f"image_processing_worker::on_message A new video arrived, starting to process it. Its timestamp going to be: {date_time}")

        copy_results = copy_video_from_edge(source_ip="192.168.178.64")

        if not copy_results:
            raise ValueError("Raw video has not been copied from the edge. This is a fatal error.")
        
        video_processing_outcome = process_video(src_file=video_source_file, timestamp=date_time)
        
        print("success:", video_processing_outcome)

        #let's convert it to the right format now 
        location = f'/opt/ready_videos/{date_time}.mp4'
        ffmpeg_command = [
            'ffmpeg',
            '-i', '/opt/processed_videos/processed.mp4',
            '-profile:v', 'main',
            '-movflags', '+faststart',
            location
        ]

        # Run the ffmpeg command
        try:
            subprocess.run(ffmpeg_command, check=True)
            print("Conversion completed successfully!")
        except subprocess.CalledProcessError as e:
            print(f"Error: {e}")
                        
        if video_processing_outcome:
            client.publish(topic, "new_processed_video")

def copy_video_from_edge(source_ip):
    # Password for edge layer, this could be moved to env variable 
    password = 'raspberry'
    source = f"cm3070@{source_ip}:/opt/unprocessed_videos/raw.mp4"
    destination = '/opt/unprocessed_videos'

    # SSHpass command to allow using password
    sshpass_command = f"sshpass -p {password} scp -r {source} {destination}"

    try:
        subprocess.run(sshpass_command, shell=True, check=True)
        print("Video copied from edge successfully!")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error: {e}")
        return False 

def process_video(src_file,timestamp):
    video=cv2.VideoCapture(src_file)
    # Check if the video opened successfully
    if (video.isOpened()== False): 
        print("Error opening video file")

    initial_frame = None

    # Define the codec and create VideoWriter object
    fourcc = cv2.VideoWriter_fourcc(*'mp4v')

    print(f"image_processing_worker::process_video, Its timestamp going to be: {timestamp} ")
    out = cv2.VideoWriter(f"{destination_dir}/processed.mp4", fourcc, 20.0,  (int(video.get(3)), int(video.get(4))))
    while video.isOpened():
        status=0
        # Capture frame by frame
        check, frame = video.read()
        if not check:
            print("not check, break")
            break
        if frame is None:
            print("frame is none")
            break 
        time.sleep(0.01)

        # Gray conversion and noise reduction (smoothening)
        gray_frame=cv2.cvtColor(frame[:,:,:],cv2.COLOR_BGR2GRAY)
        blur_frame=cv2.GaussianBlur(gray_frame,(25,25),0)

        # First captured frame is the baseline image
        if initial_frame is None:
            initial_frame = blur_frame
            continue
        
        # let's get the difference between baseline and new frame
        delta_frame=cv2.absdiff(initial_frame,blur_frame)
        #  here we convert delta_frame into binary image
        threshold_frame=cv2.threshold(delta_frame,30,255, cv2.THRESH_BINARY)[1]

        # identify all the contours in the image.
        (contours,_)=cv2.findContours(threshold_frame,cv2.RETR_EXTERNAL,cv2.CHAIN_APPROX_SIMPLE)

        for c in contours:
            #  filter out small contours
            if cv2.contourArea(c) < 1000:
                continue
            status=status + 1
            (x, y, w, h)=cv2.boundingRect(c)
            cv2.rectangle(frame, (x, y), (x+w, y+h), (0,255,0), 1)
        
        # write the flipped frame
        out.write(frame)

        # allow stopping app with 'q'    
        if cv2.waitKey(1) == ord('q'):
            break

    # release video object
    out.release()
    video.release()

    # destroy windows
    cv2.destroyAllWindows()

    return True 

def main():
    client = mqtt.Client()
    client.on_message = on_message
    client.connect(mqtt_address)
    client.subscribe(topic)
    client.loop_forever()
 

if __name__ == "__main__":
    main()