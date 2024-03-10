# Importing the Required Library
from twilio.rest import Client
import os 


class SmsWorker:
    def __init__(self) -> None:
        pass

    def send_sms(self,message_content,from_number,to_number):
        account_sid = os.environ.get('twilio_account_sid')
        auth_token = os.environ.get('twilio_auth_token')

        client = Client(account_sid, auth_token)

        # Sending an SMS
        message = client.messages.create(
            body=message_content,
            from_=from_number,  # Your Twilio phone number
            to=to_number  # Recipient's phone number
        )

        # Verifying if the message is sent or not
        print(f"SMS sent successfully. Message ID: {message.sid}")