# Crispcam Demo Outline

There are two flows for this demo:

1. Anthos Service Mesh
2. Open Source

You can show both and interleave them - but they have different flavours.

## Step 1 - Login to the cluster

```
gcloud beta container clusters get-credentials demo-cluster-a --region europe-west1 --project crisp-retail-demo
```

## Step 2 - Load Kiali and Jaeger Proxies (Open Source Only)

**Note: You only need to do this if demonstrating open source**

```bash
trap 'kill %1; kill %2' SIGINT; \
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=kiali -o jsonpath='{.items[0].metadata.name}') 20001:20001 | sed -e 's/^/[Kiali ] /' & \
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=jaeger -o jsonpath='{.items[0].metadata.name}') 15032:16686 | sed -e 's/^/[Jaeger] /' & \
wait
trap - SIGINT
```

## Step 3 - Demo the app

Visit [crisps.gcp-north.co.uk](https://crisps.gcp-north.co.uk) and show the various features (reviews, search, etc) and explain how it's all microservices.

Show [Crispcam](https://crisps.gcp-north.co.uk/crispcam) and talk a bit about Vision AutoML

## Step 4 - Demo Service Visibility

### Google Cloud

Visit the Anthos Service Mesh Services page in the cloud console to see the various graph, timeline and topology views.

Show SLO and SLI metrics at this stage

### Open Source

Visit the [kiali dashboard](http://localhost:20001/kiali/) (admin/admin) and click Graph - ensure the default project is selected (it's not by default).

Show how this is all handled without any major code changes (just some http headers). If the customer wants more detail the attached file [code/RestTemplateConfig.kt](code/RestTemplateConfig.kt) shows how a Spring Boot config bean could be used to achieve this and lists the http headers persisted.

Show how load balancing (e.g. to the review app) is handled and if you add the traffic percentage label show how it is roughly 50/50.

## Step 5 - Show client-side load balancing

### 5a - Background

Explain how normally a developer would need to write code to get client-side load balancing working. This might be a library and depending on the language it adds overhead.

Also what is the LB strategy an app should use? Isn't this an infrastructure problem?

Open the file [yaml/virtual-service-lb.yaml](yaml/virtual-service-lb.yaml) and show how it's currently doing that for you - load balancing 50/50 between red and yellow. On the app itself, if you refresh the page a few times you'll see it changing too.

### 5b - Firefox and Chrome

Open the file [yaml/virtual-service-firefox.yaml](yaml/virtual-service-firefox.yaml) and show how you've made a simple rule to look at the http headers - in this case the user agent. Explain that this is for a backend (not frontend) service - but because we persist the headers we can use this as part of our client-side load balancing.

Apply the config:

```
kubectl apply -f yaml/virtual-service-firefox.yaml
```

Show how Firefox only shows red stars now and Chrome only shows yellow

## Step 6 - Show circuit breaking

### 6a - Background

Point out the 'broken' service - it holds a TCP connection for 10 seconds before returning - by default most software will just wait for it to return.

```kotlin
Thread.sleep(10000)
```

Introduce circuit breaking - how a broken service shouldn't break your app and instead should be marked dead (a blackout is better than a brownout). Explain how it might be a database, storage or other problems. In many cases you can do this in software (similar to LB) but it's complex and language-specific. Sharing that a broken service is down is complex too - you need to use something like Netflix's Hystrix and manage/configure that.

### 6b - Break the service

Show the file [yaml/virtual-service-broken.yaml](yaml/virtual-service-broken.yaml) and how it routes 100% of traffic to this broken service. Then apply it:

```
kubectl apply -f yaml/virtual-service-broken.yaml
```

Go back to [crisps.gcp-north.co.uk](https://crisps.gcp-north.co.uk) and show it's broken!

### 6c - Fix it

Show the [yaml/destination-rule.yaml](yaml/destination-rule.yaml) file and show the circuit breaker configuration.

Now show [yaml/virtual-service-fixed.yaml](yaml/virtual-service-fixed.yaml) - point out how it now has a `timeout: 1s` parameter - if the request takes longer than that, intercept and return a gateway error.

Apply this change:

```
kubectl apply -f yaml/virtual-service-fixed.yaml
```

Now go back to [crisps.gcp-north.co.uk](https://crisps.gcp-north.co.uk) and observe the reviews are showing down, but the site is working properly.

### 6d - Show in Kiali

Go back to the [kiali dashboard](http://localhost:20001/kiali/) and show how the service is now going from 'healthy' to 'unhealthy' - and we didn't change any code!

## Step 7 - Show tracing

Show Jaeger tracing - [Jaeger Dashboard](http://localhost:15032/)

## Step 9 - Cleanup

Apply the cleanup file so it works for everyone again:

```
kubectl apply -f yaml/virtual-service-lb.yaml
```
