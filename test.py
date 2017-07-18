import requests

resp = requests.post('http://localhost:8001/api/v1/queue/test', data={"test":True})
print(resp.text)

resp = requests.get('http://localhost:8001/api/v1/queue/test')
print(resp.text)

resp = requests.post('http://localhost:8001/api/v1/queue/test', json={"test":True})
print(resp.text)
