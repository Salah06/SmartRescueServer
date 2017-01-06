import requests

for i in range(0, 1000) :
	for j in range(0, 99) :
		r = requests.post("http://localhost:1234/androidEmergency", data={'id': str(i)})

print "finish"