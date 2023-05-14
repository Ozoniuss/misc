import random

with open("a.txt", "w") as f:
    for i in range(10000):
        n = random.randint(0, 255)
        f.write(str(n))
        f.write(",")