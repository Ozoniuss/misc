import datetime
import json
import socket
import sys
from collections import Counter
from matplotlib import pyplot as plt
import numpy as np
import pandas as pd
from dataclasses import dataclass, asdict


HOST = "127.0.0.1"  
PORT = sys.argv[1]

if PORT == None or PORT == "":
    PORT = 13311
else:
    PORT = int(PORT)

@dataclass
class ArticleLikedEvent:
    event_id: str
    article_id: str
    timestamp: str

@dataclass
class Message:
    id: int
    events: list[ArticleLikedEvent]

@dataclass
class Ack: 
    id: int

def parse_time(timestamp: str) -> datetime.datetime:
    """Helper to generate a valid timestamp from a string"""
    
    # Parse the string into a datetime object
    dt = datetime.datetime.fromisoformat(timestamp[:-6])

    # Extract timezone offset
    tz_offset = timestamp[-6:]

    # Parse the timezone offset
    hours_offset = int(tz_offset[:3])
    minutes_offset = int(tz_offset[4:])

    # Create a timedelta object for timezone offset
    td = datetime.timedelta(hours=hours_offset, minutes=minutes_offset)

    # Adjust datetime object with timezone offset
    if tz_offset[3] == '+':
        dt = dt - td
    else:
        dt = dt + td

    return dt

def generate_like_graph(events: list[ArticleLikedEvent]):
    
    totals = []
    for e in events:
        ts = parse_time(e.timestamp)
        totals.append(datetime.datetime(
        ts.year,
        ts.month,
        ts.day,
        ts.hour,
        ts.minute,
        ts.second))

    totalsDict = dict(Counter(totals))
    xaxis, yaxis = [], []
    for k in sorted(totalsDict.keys()):
        xaxis.append(k)
        yaxis.append(totalsDict[k])

    #define data
    df = pd.DataFrame({
        'time':[str(t.second) for t in xaxis],
        'seconds': yaxis})

    #plot time series
    plt.plot(df.time, df.seconds, linewidth=3)
    #add title and axis labels
    plt.title('Likes by second')
    plt.xlabel('Second (trimmed)')
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

            # generate the graph upon ctrl + C
            except KeyboardInterrupt:
                generate_like_graph(receivedEvents)
                break
            
            except Exception as e:
                print("exception", e)
                break

