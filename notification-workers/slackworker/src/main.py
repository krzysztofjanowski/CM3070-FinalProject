import requests

class SlackWorker:
    
    def __init__(self) -> None:
        pass

    def send_status_to_slack(self,api_token="", message="test"):

        slack_data = {
                    "blocks": [
                        {
                            "type": "header",
                            "text": {
                                "type": "plain_text",
                                "text": f"{message}"
                            }
                        },
                        {
                            "type": "section",
                            "text": {
                                "type": "mrkdwn",
                                "text": f"*{message}*"
                            }
                        },
                    ]
                }
        
        requests.post(f"https://hooks.slack.com/services/{api_token}", json=slack_data)

