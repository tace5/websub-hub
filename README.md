To notify subscribers of a topic about an update, send a post request to localhost:8080/notify:
```
curl -X POST -d 'hub.topic=/a/topic' -d 'data=Hello World!' localhost:8080/notify
```
