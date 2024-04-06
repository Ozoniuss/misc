import datetime
import json
import socket
from dataclasses import dataclass, asdict

HOST = "127.0.0.1"  
PORT = 13311

@dataclass
class ArticleLikedEvent:
    event_id: str
    article_id: str
    timestamp: datetime.datetime

@dataclass
class Message:
    id: int
    events: list[ArticleLikedEvent]

@dataclass
class Ack: 
    id: int

receivedEvents: list[ArticleLikedEvent] = []


with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((HOST, PORT))
    s.listen()
    conn, addr = s.accept()
    with conn:
        print(f"Connected by {addr}")
        while True:
            try:
                data = conn.recv(4096).decode('UTF-8')
                message = Message(**json.loads(str(data)))
                ack = json.dumps(asdict(Ack(id=message.id)))
                receivedEvents.extend(message.events)
                conn.sendall(bytes(ack,encoding="utf-8"))

            except Exception as e:
                print("exception", e)
                break