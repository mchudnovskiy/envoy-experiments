# Envoy's Active Health Check discovery
## Description

During this task we will test envoy active health check on cluster with two http endpoins. 

Cluster configuration:
```yaml
    load_assignment:
      cluster_name: external_service
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: first
                port_value: 7777
        - endpoint:
            address:
              socket_address:
                address: second
                port_value: 7777  
```
Health check configuration:
```yaml
    health_checks:
      timeout: 2s
      interval: 1s
      unhealthy_threshold: 3
      healthy_threshold: 3
      no_traffic_interval: 60s
      event_log_path: /dev/stdout
      always_log_health_check_failures: false
      http_health_check:
        path: /  
```


Full envoy configuration and endpoint sources are avaliable on [Github][envoy-exp-github]

All components are running locally in docker with the sharing docker bridge network:
```bash
docker run --rm --name=first -it --net=envoy-exp  -d simple_http

docker run --rm --name=second -it --net=envoy-exp  -d simple_http

docker run --rm --name=envoy --net=envoy-exp -it -d -v=$PWD:/config  -p 8888:8888 -p 9999:9999 envoyproxy/envoy /usr/local/bin/envoy -c config/proxy.yaml
```
Test case actions: 
1. All endpoints are healthy. Run envoy and two http endpoints locally in docker and check logs and envoy `/clusters` output.
2. The first endpoint becomes unhealthy. Stop the first enpoint and check  logs and envoy `/clusters`  output.
3. All endpoints are healthy again. Start the first endpoint and check  logs and envoy `/clusters`  output.

Expected results:
  1. All endpoints are healthy:
      -  envoy `/clusters` output has two healthy enpoints, and all endpoints are getting envoy's hc-requests at one min period without traffic or one sec period with traffic
  2. One unhealthy endpoint:
      - envoy `/clusters` output has one 'pending_dynamic_removal' endpoint that will be removed after 3 unsuccessful hc-requests
  3. All endpoints are healthy again.
      -  envoy `/clusters` output has two healthy enpoints after 3 successfull hc-requests to the first endpoint, and all endpoints are getting envoy's hc-requests

## All endpoints are healthy
After start on both endpoints we have envoy’s hc-requests with one min period (no_traffic_interval: 60s).

Endpoints log sample:
```
 2019-08-05 12:28:11.9706493 +0000 UTC m=+524.182744901 Got request.
Method: GET
 URL: /
 Header: map[Content-Length:[0] User-Agent:[Envoy/HC]]
 Body: {}


 2019-08-05 12:29:11.9045269 +0000 UTC m=+584.184212401 Got request.
Method: GET
 URL: /
 Header: map[Content-Length:[0] User-Agent:[Envoy/HC]]
 Body: {}
```
Envoy log has healthy_event records for both endpoints:
```
{"health_checker_type":"HTTP","host":{"socket_address":{"protocol":"TCP","address":"172.18.0.3","resolver_name":"","ipv4_compat":false,"port_value":7777}},"cluster_name":"external_service","add_healthy_event":{"first_check":true},"timestamp":"2019-08-05T12:28:11.971Z"}

{"health_checker_type":"HTTP","host":{"socket_address":{"protocol":"TCP","address":"172.18.0.2","resolver_name":"","ipv4_compat":false,"port_value":7777}},"cluster_name":"external_service","add_healthy_event":{"first_check":true},"timestamp":"2019-08-05T12:28:11.972Z"}
```

*Notice: The first and second endpoints received addresses 172.18.0.2 and 172.18.0.3 respectively.*

After a few requests to envoy (`curl localhost:8888/test`) we can see, that  hc period has been changed to one sec.

Endpoints log sample:
```
2019-08-05 12:54:34.4712572 +0000 UTC m=+1750.500992101 Got request.
Method: GET
 URL: /
 Header: map[Content-Length:[0] User-Agent:[Envoy/HC]]
 Body: {}


 2019-08-05 12:54:35.4674201 +0000 UTC m=+1751.497154901 Got request.
Method: GET
 URL: /
 Header: map[Content-Length:[0] User-Agent:[Envoy/HC]]
```

All requests are balancing between endpoints with round robin strategy.

Envoy `/clusters` output: 
```
external_service::default_priority::max_connections::1024
external_service::default_priority::max_pending_requests::1024
external_service::default_priority::max_requests::1024
external_service::default_priority::max_retries::3
external_service::high_priority::max_connections::1024
external_service::high_priority::max_pending_requests::1024
external_service::high_priority::max_requests::1024
external_service::high_priority::max_retries::3
external_service::added_via_api::false
external_service::172.18.0.2:7777::cx_active::4
external_service::172.18.0.2:7777::cx_connect_fail::0
external_service::172.18.0.2:7777::cx_total::4
external_service::172.18.0.2:7777::rq_active::0
external_service::172.18.0.2:7777::rq_error::0
external_service::172.18.0.2:7777::rq_success::21
external_service::172.18.0.2:7777::rq_timeout::0
external_service::172.18.0.2:7777::rq_total::21
external_service::172.18.0.2:7777::hostname::first
external_service::172.18.0.2:7777::health_flags::healthy
external_service::172.18.0.2:7777::weight::1
external_service::172.18.0.2:7777::region::
external_service::172.18.0.2:7777::zone::
external_service::172.18.0.2:7777::sub_zone::
external_service::172.18.0.2:7777::canary::false
external_service::172.18.0.2:7777::priority::0
external_service::172.18.0.2:7777::success_rate::-1
external_service::172.18.0.2:7777::local_origin_success_rate::-1
external_service::172.18.0.3:7777::cx_active::4
external_service::172.18.0.3:7777::cx_connect_fail::0
external_service::172.18.0.3:7777::cx_total::4
external_service::172.18.0.3:7777::rq_active::0
external_service::172.18.0.3:7777::rq_error::0
external_service::172.18.0.3:7777::rq_success::20
external_service::172.18.0.3:7777::rq_timeout::0
external_service::172.18.0.3:7777::rq_total::20
external_service::172.18.0.3:7777::hostname::second
external_service::172.18.0.3:7777::health_flags::healthy
external_service::172.18.0.3:7777::weight::1
external_service::172.18.0.3:7777::region::
external_service::172.18.0.3:7777::zone::
external_service::172.18.0.3:7777::sub_zone::
external_service::172.18.0.3:7777::canary::false
external_service::172.18.0.3:7777::priority::0
external_service::172.18.0.3:7777::success_rate::-1
external_service::172.18.0.3:7777::local_origin_success_rate::-1
```
*Notice: we have 2  endpoints (172.18.0.2,172.18.0.3) in cluster external_service and they are healthy (`health_flags::healthy`)*


## One unhealthy endpoint
After we turn off the first endpoint, the envoy  starts the procedure of endpoint removal from the cluster after the first unsuccessful HC-request. 

This can be seen on the `health_flags::/pending_dynamic_removal` for the first endpoint at 172.18.0.2 

Envoy `/clusters` output: 
```
external_service::default_priority::max_connections::1024
external_service::default_priority::max_pending_requests::1024
external_service::default_priority::max_requests::1024
external_service::default_priority::max_retries::3
external_service::high_priority::max_connections::1024
external_service::high_priority::max_pending_requests::1024
external_service::high_priority::max_requests::1024
external_service::high_priority::max_retries::3
external_service::added_via_api::false
external_service::172.18.0.2:7777::cx_active::0
external_service::172.18.0.2:7777::cx_connect_fail::0
external_service::172.18.0.2:7777::cx_total::4
external_service::172.18.0.2:7777::rq_active::0
external_service::172.18.0.2:7777::rq_error::0
external_service::172.18.0.2:7777::rq_success::21
external_service::172.18.0.2:7777::rq_timeout::0
external_service::172.18.0.2:7777::rq_total::21
external_service::172.18.0.2:7777::hostname::first
external_service::172.18.0.2:7777::health_flags::/pending_dynamic_removal
external_service::172.18.0.2:7777::weight::1
external_service::172.18.0.2:7777::region::
external_service::172.18.0.2:7777::zone::
external_service::172.18.0.2:7777::sub_zone::
external_service::172.18.0.2:7777::canary::false
external_service::172.18.0.2:7777::priority::0
external_service::172.18.0.2:7777::success_rate::-1
external_service::172.18.0.2:7777::local_origin_success_rate::-1
external_service::172.18.0.3:7777::cx_active::4
external_service::172.18.0.3:7777::cx_connect_fail::0
external_service::172.18.0.3:7777::cx_total::4
external_service::172.18.0.3:7777::rq_active::0
external_service::172.18.0.3:7777::rq_error::0
external_service::172.18.0.3:7777::rq_success::20
external_service::172.18.0.3:7777::rq_timeout::0
external_service::172.18.0.3:7777::rq_total::20
external_service::172.18.0.3:7777::hostname::second
external_service::172.18.0.3:7777::health_flags::healthy
external_service::172.18.0.3:7777::weight::1
external_service::172.18.0.3:7777::region::
external_service::172.18.0.3:7777::zone::
external_service::172.18.0.3:7777::sub_zone::
external_service::172.18.0.3:7777::canary::false
external_service::172.18.0.3:7777::priority::0
external_service::172.18.0.3:7777::success_rate::-1
external_service::172.18.0.3:7777::local_origin_success_rate::-1
```

After three unsuccessful hc-requests envoy marks endpoint 172.18.0.2 as 'unhealthy' and removes it from cluster external_service.

Envoy log:
```
{"health_checker_type":"HTTP","host":{"socket_address":{"protocol":"TCP","address":"172.18.0.2","resolver_name":"","ipv4_compat":false,"port_value":7777}},"cluster_name":"external_service","eject_unhealthy_event":{"failure_type":"NETWORK"},"timestamp":"2019-08-05T13:05:21.471Z"}
```
Envoy `/clusters` output after 3rd unsuccessful hc:
```external_service::default_priority::max_connections::1024
external_service::default_priority::max_pending_requests::1024
external_service::default_priority::max_requests::1024
external_service::default_priority::max_retries::3
external_service::high_priority::max_connections::1024
external_service::high_priority::max_pending_requests::1024
external_service::high_priority::max_requests::1024
external_service::high_priority::max_retries::3
external_service::added_via_api::false
external_service::172.18.0.3:7777::cx_active::4
external_service::172.18.0.3:7777::cx_connect_fail::0
external_service::172.18.0.3:7777::cx_total::4
external_service::172.18.0.3:7777::rq_active::0
external_service::172.18.0.3:7777::rq_error::0
external_service::172.18.0.3:7777::rq_success::20
external_service::172.18.0.3:7777::rq_timeout::0
external_service::172.18.0.3:7777::rq_total::20
external_service::172.18.0.3:7777::hostname::second
external_service::172.18.0.3:7777::health_flags::healthy
external_service::172.18.0.3:7777::weight::1
external_service::172.18.0.3:7777::region::
external_service::172.18.0.3:7777::zone::
external_service::172.18.0.3:7777::sub_zone::
external_service::172.18.0.3:7777::canary::false
external_service::172.18.0.3:7777::priority::0
external_service::172.18.0.3:7777::success_rate::-1
external_service::172.18.0.3:7777::local_origin_success_rate::-1
```

*Notice that only one endpoint 172.18.0.3 remained in the cluster external_service.
All requests (`curl localhost:8888/test`) are routing to healthy endpoint 172.18.0.3* 


## All endpoints are healthy again

When disabled endpoint starts again envoy adds it to cluster after 3rd successful HC-request.

Envoy log
```
{"health_checker_type":"HTTP","host":{"socket_address":{"protocol":"TCP","address":"172.18.0.2","resolver_name":"","ipv4_compat":false,"port_value":7777}},"cluster_name":"external_service","add_healthy_event":{"first_check":true},"timestamp":"2019-08-05T13:13:16.974Z"}
```

Envoy `/clusters` output 
```
external_service::default_priority::max_connections::1024
external_service::default_priority::max_pending_requests::1024
external_service::default_priority::max_requests::1024
external_service::default_priority::max_retries::3
external_service::high_priority::max_connections::1024
external_service::high_priority::max_pending_requests::1024
external_service::high_priority::max_requests::1024
external_service::high_priority::max_retries::3
external_service::added_via_api::false
external_service::172.18.0.2:7777::cx_active::0
external_service::172.18.0.2:7777::cx_connect_fail::0
external_service::172.18.0.2:7777::cx_total::0
external_service::172.18.0.2:7777::rq_active::0
external_service::172.18.0.2:7777::rq_error::0
external_service::172.18.0.2:7777::rq_success::0
external_service::172.18.0.2:7777::rq_timeout::0
external_service::172.18.0.2:7777::rq_total::0
external_service::172.18.0.2:7777::hostname::first
external_service::172.18.0.2:7777::health_flags::healthy
external_service::172.18.0.2:7777::weight::1
external_service::172.18.0.2:7777::region::
external_service::172.18.0.2:7777::zone::
external_service::172.18.0.2:7777::sub_zone::
external_service::172.18.0.2:7777::canary::false
external_service::172.18.0.2:7777::priority::0
external_service::172.18.0.2:7777::success_rate::-1
external_service::172.18.0.2:7777::local_origin_success_rate::-1
external_service::172.18.0.3:7777::cx_active::4
external_service::172.18.0.3:7777::cx_connect_fail::0
external_service::172.18.0.3:7777::cx_total::4
external_service::172.18.0.3:7777::rq_active::0
external_service::172.18.0.3:7777::rq_error::0
external_service::172.18.0.3:7777::rq_success::20
external_service::172.18.0.3:7777::rq_timeout::0
external_service::172.18.0.3:7777::rq_total::20
external_service::172.18.0.3:7777::hostname::second
external_service::172.18.0.3:7777::health_flags::healthy
external_service::172.18.0.3:7777::weight::1
external_service::172.18.0.3:7777::region::
external_service::172.18.0.3:7777::zone::
external_service::172.18.0.3:7777::sub_zone::
external_service::172.18.0.3:7777::canary::false
external_service::172.18.0.3:7777::priority::0
external_service::172.18.0.3:7777::success_rate::-1
external_service::172.18.0.3:7777::local_origin_success_rate::-1
```

*Notice: there are two endpoints in cluster again, 172.18.0.2 and 172.18.0.3. All endpoints have `health_flags::healthy`*

All endpoints receive envoy hc-requests.

Endpoints log sample:
```
2019-08-05 13:25:07.3861482 +0000 UTC m=+712.030993701 Got request.
Method: GET
 URL: /
 Header: map[Content-Length:[0] User-Agent:[Envoy/HC]]
 Body: {}

 2019-08-05 13:25:08.3499997 +0000 UTC m=+713.029842801 Got request.
Method: GET
 URL: /
 Header: map[Content-Length:[0] User-Agent:[Envoy/HC]]
 Body: {}
 ```

## Conclusion
After all planned actions we've got expected results. Envoy active health check is usefull thing for external services, that runs outside k8s clusters.
 
[envoy-exp-github]: https://github.com/mchudnovskiy/envoy-experiments

