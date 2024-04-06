import datetime
import json
import socket
from collections import Counter
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

PORT = 13312

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

def generate_like_graph(events: list[ArticleLikedEvent]):
    
    totals = []
    for e in events:
        totals.append(str(datetime.datetime(
        e.timestamp.year,
        e.timestamp.month,
        e.timestamp.day,
        e.timestamp.hour,
        e.timestamp.minute,
        e.timestamp.second)))

    totalsDict = dict(Counter(totals))
    xaxis, yaxis = [], []
    for k in sorted(totalsDict.keys()):
        xaxis.append(k)
        yaxis.append(totalsDict[k])

    #define data
    df = pd.DataFrame({
        'time':xaxis,
        'seconds': yaxis})

    #plot time series
    plt.plot(df.time, df.seconds, linewidth=3)
    #add title and axis labels
    plt.title('Likes by second')
    plt.xlabel('Second')
    plt.ylabel('Likes')

    plt.savefig('data.png')

receivedEvents: list[ArticleLikedEvent] = []


with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((HOST, PORT))
    s.listen()
    conn, addr = s.accept()
    c = 0
    with conn:
        print(f"Connected by {addr}")
        while True:
            try:
                c += 1
                data = conn.recv(1024 * 1024 * 10).decode('UTF-8')
                print("received data", data)
                message = Message(**json.loads(str(data)))
                message.events = [ArticleLikedEvent(**e) for e in message.events]

                receivedEvents.extend(message.events)


                ack = json.dumps(asdict(Ack(id=message.id)))
                conn.sendall(bytes(ack,encoding="utf-8"))
                all = [e.event_id for e in receivedEvents]
                print("allevents", all, len(all))

                if c == 10:

                    break
            except KeyboardInterrupt:
                generate_like_graph(receivedEvents)
            
            except Exception as e:
                print("exception", e)
                break

