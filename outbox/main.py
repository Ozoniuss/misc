import datetime
import json
import socket
from time import sleep
from typing import Counter
from matplotlib import pyplot as plt
import numpy as np
import pandas as pd
from dataclasses import dataclass, asdict


# #define data
# df = pd.DataFrame({'date': np.array([datetime.datetime(2020, 1, i+1)
#                                      for i in range(12)]),
#                    'sales': [3, 4, 4, 7, 8, 9, 14, 17, 12, 8, 8, 13]})

# #plot time series
# plt.plot(df.date, df.sales, linewidth=3)
# plt.savefig("data.png")

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
                data = conn.recv(1024 * 1024 * 10).decode('UTF-8')
                print("received data", data)
                message = Message(**json.loads(str(data)))
                receivedEvents.extend(message.events)


                ack = json.dumps(asdict(Ack(id=message.id)))
                conn.sendall(bytes(ack,encoding="utf-8"))
                all = [e['event_id'] for e in receivedEvents]
                print("allevents", all, len(all))

            except Exception as e:
                print("exception", e)
                break

