import os
import yagmail


class EmailWorker:
    def __init__(self) -> None:
        pass


    def send_email(self,subject="A new video recorded", contents="A new video has been uploaded to the dashboard",dst_email=""):
        """
        Sends an email with the subject and contents 
        :param subject: subject of the email
        :param contents: email contents
        :param dst_email: email address to which emails are sent to
        :return:
        """

        # from email is fixed but could be refactored
        yag = yagmail.SMTP('nobigmic@gmail.com',os.environ.get('gmail_pass') )

        yag.send(to=dst_email,subject=subject, contents=contents)
