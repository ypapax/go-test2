# Task
This test is intended to test your ability to write a complete micro-service backend. This service should be capable of serving the dataset you downloaded from NERACOOS via a RESTful API. For this test you can use either Python or GO and a database of your choice (e.g., MySQL, Postgres, etc.). Additionally, the service should be accessible at [your.domain.com]/test/api/v1. If you do not have a domain name an IP address will be sufficient. All data should be served as json containing appropriate results when available and an accompanying error message when appropriate. The json attribute names are up to you.

Your service should include the following API endpoints:

1) current_speed,
2) temperature,

All endpoints should be able to filter results based on the presence of the parameters start, stop, min, max and avg. For example, the endpoint your.domain.com/test/api/v1/current_speed?start=01/01/2016 should return all current speed records since 01/01/2016 whereas the endpoint your.domain.com/test/api/v1/current_speed/min?start=01/01/2016 should return the min current speed since 01/01/2016. Be sure to check those _qc values :-).

Your service should handle and respond to reasonable errors accordingly. By reasonable we mean dates that are out of bounds or passing parameters that are not supported (e.g., your.domain.com/api/v1/current_speed/unsupported_foo?stop=01/01/2016). Of course, your service should provide appropriate responses (e.g., 404) for endpoints that do not exist.
# Solution

## endpoints
```
curl -sS localhost:8181/test/api/v1/temperature/max?start=01/01/2016
{"result":19.601141}

curl -sS localhost:8181/test/api/v1/current_speed/avg
{"result":21.631215155609603}

curl -sS localhost:8181/test/api/v1/temperature/min?start=01/01/2016&stop=01/01/2017
{"result":6.325187}

curl -sS localhost:8181/test/api/v1/fake_ep/fake?stop=01/01/2006
not supported endpoint

curl -sS localhost:8181/test/api/v1/current_speed?start=01/01/2016&stop=05/12/2016
[{"current_speed":19.3578,"time":"2016-12-05T03:00:00+03:00"},{"current_speed":22.8774,"time":"2016-12-05T02:40:00+03:00"},{"current_speed":24.6372,"time":"2016-12-05T02:20:00+03:00"},{"current_speed"...

curl -sS localhost:8181/test/api/v1/current_speed?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}
```


## deploy
```
$GOPATH/src/github.com/ypapax/go-test2/deploy.sh user host port mongoConnectionString
```
for example:
```
$GOPATH/src/github.com/ypapax/go-test2/deploy.sh user1 178.116.19.128 8822 10.133.131.30:27112,10.133.36.252:27112,10.133.7.151:27112,10.133.7.90:27112
```
will deploy the service on `178.116.19.128` port `8822` connectin via ssh as `user1` during deployment.
## check
`./check.sh` requests to different endpoints:
```
cd $GOPATH/github.com/ypapax/go-test2 && ./check.sh localhost:8181
```
```
curl -sS localhost:8181/test/api/v1/temperature/max?start=01/01/2016
{"result":19.601141}
curl -sS localhost:8181/test/api/v1/temperature/max
{"result":19.601141}
curl -sS localhost:8181/test/api/v1/temperature/max?start=01/01/2016&stop=01/01/2017
{"result":19.601141}
curl -sS localhost:8181/test/api/v1/temperature/max?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}

curl -sS localhost:8181/test/api/v1/temperature/avg?start=01/01/2016
{"result":12.747452368335534}
curl -sS localhost:8181/test/api/v1/temperature/avg
{"result":12.697609162955457}
curl -sS localhost:8181/test/api/v1/temperature/avg?start=01/01/2016&stop=01/01/2017
{"result":12.747452368335534}
curl -sS localhost:8181/test/api/v1/temperature/avg?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}

curl -sS localhost:8181/test/api/v1/temperature/fake?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/temperature/fake
not supported endpoint

curl -sS localhost:8181/test/api/v1/temperature/fake?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/temperature/fake?stop=01/01/2006
not supported endpoint


curl -sS localhost:8181/test/api/v1/temperature/min?start=01/01/2016
{"result":6.325187}
curl -sS localhost:8181/test/api/v1/temperature/min
{"result":6.325187}
curl -sS localhost:8181/test/api/v1/temperature/min?start=01/01/2016&stop=01/01/2017
{"result":6.325187}
curl -sS localhost:8181/test/api/v1/temperature/min?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}

curl -sS localhost:8181/test/api/v1/temperature?start=01/01/2016
[{"temperature":8.3420925,"time":"2016-12-05T17:00:00+03:00"},{"temperature":8.3420925,"time":"2016-

curl -sS localhost:8181/test/api/v1/temperature
[{"temperature":8.3420925,"time":"2016-12-05T17:00:00+03:00"},{"temperature":8.3420925,"time":"2016-

curl -sS localhost:8181/test/api/v1/temperature?start=01/01/2016&stop=01/01/2017
[{"temperature":8.3420925,"time":"2016-12-05T17:00:00+03:00"},{"temperature":8.3420925,"time":"2016-

curl -sS localhost:8181/test/api/v1/temperature?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}



curl -sS localhost:8181/test/api/v1/fake_ep/max?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/max
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/max?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/max?stop=01/01/2006
not supported endpoint


curl -sS localhost:8181/test/api/v1/fake_ep/avg?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/avg
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/avg?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/avg?stop=01/01/2006
not supported endpoint


curl -sS localhost:8181/test/api/v1/fake_ep/fake?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/fake
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/fake?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/fake?stop=01/01/2006
not supported endpoint


curl -sS localhost:8181/test/api/v1/fake_ep/min?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/min
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/min?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep/min?stop=01/01/2006
not supported endpoint


curl -sS localhost:8181/test/api/v1/fake_ep?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/fake_ep?stop=01/01/2006
not supported endpoint




curl -sS localhost:8181/test/api/v1/current_speed/max?start=01/01/2016
{"result":66.8724}
curl -sS localhost:8181/test/api/v1/current_speed/max
{"result":81.2441}
curl -sS localhost:8181/test/api/v1/current_speed/max?start=01/01/2016&stop=01/01/2017
{"result":66.8724}
curl -sS localhost:8181/test/api/v1/current_speed/max?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}

curl -sS localhost:8181/test/api/v1/current_speed/avg?start=01/01/2016
{"result":18.003292556142668}
curl -sS localhost:8181/test/api/v1/current_speed/avg
{"result":21.631215155609603}
curl -sS localhost:8181/test/api/v1/current_speed/avg?start=01/01/2016&stop=01/01/2017
{"result":18.003292556142668}
curl -sS localhost:8181/test/api/v1/current_speed/avg?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}

curl -sS localhost:8181/test/api/v1/current_speed/fake?start=01/01/2016
not supported endpoint

curl -sS localhost:8181/test/api/v1/current_speed/fake
not supported endpoint

curl -sS localhost:8181/test/api/v1/current_speed/fake?start=01/01/2016&stop=01/01/2017
not supported endpoint

curl -sS localhost:8181/test/api/v1/current_speed/fake?stop=01/01/2006
not supported endpoint


curl -sS localhost:8181/test/api/v1/current_speed/min?start=01/01/2016
{"result":0}
curl -sS localhost:8181/test/api/v1/current_speed/min
{"result":0}
curl -sS localhost:8181/test/api/v1/current_speed/min?start=01/01/2016&stop=01/01/2017
{"result":0}
curl -sS localhost:8181/test/api/v1/current_speed/min?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}

curl -sS localhost:8181/test/api/v1/current_speed?start=01/01/2016
[{"current_speed":6.4526,"time":"2016-12-05T17:00:00+03:00"},{"current_speed":8.799,"time":"2016-12-

curl -sS localhost:8181/test/api/v1/current_speed
[{"current_speed":6.4526,"time":"2016-12-05T17:00:00+03:00"},{"current_speed":8.799,"time":"2016-12-

curl -sS localhost:8181/test/api/v1/current_speed?start=01/01/2016&stop=01/01/2017
[{"current_speed":6.4526,"time":"2016-12-05T17:00:00+03:00"},{"current_speed":8.799,"time":"2016-12-

curl -sS localhost:8181/test/api/v1/current_speed?stop=01/01/2006
{"error":"Out of bounds. Supported bounds: start=25/08/2015 stop=05/12/2016"}
```
