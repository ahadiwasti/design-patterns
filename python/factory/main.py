from abc import ABC, abstractmethod
import os

class Notification(ABC):
    @abstractmethod
    def send(self, message:str) -> str:
        pass
    
class EmailNotification(Notification):
    def __init__(self,smtp_server:str, port:int):
        self.smtp_server = smtp_server
        self.port = port
    
    def send(self, message:str)->str:
        #implement logic here
        return f"email sent via {self.smtp_server} and the message is {message}"


class WhatsappNotification(Notification):
    def __init__(self,phone_id:str,access_token:str):
        self.phone_id = phone_id
        self.access_token = access_token
    
    def send(self, message:str)->str:
        #implement the logic here
        return f"message is sent via phone id:{self.phone_id}"


def new_notification(channel:str)-> Notification:
    options = {
        "email": lambda : EmailNotification(
            smtp_server = os.getenv("SMTP SERVER Config"),
            port = 587
        ),
        "whatsapp": lambda : WhatsappNotification(
            phone_id = os.getenv("PhoneId"),
            access_token = os.getenv("accessToken")
        )
    }

    factory = options.get(channel)
    if not factory:
        raise ValueError(f"unknown channel : {channel}")
    return factory()

channels = ["email","whatsapp"]

for channel in channels:
    n = new_notification(channel)
    print(n.send("your order has shipped!!"))