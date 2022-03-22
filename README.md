# Crispcam Demo Outline

## Step 1 - Login to the cluster

```
gcloud container clusters get-credentials crisps-autopilot --region europe-west1 --project crispcam
```

## Step 2 - Demo the app

Visit [crispcam.com](https://crispcam.com) and show the various features (reviews, search, etc) and explain how it's all microservices.

Show [Crispcam](https://crispcam.com/crispcam) and talk a bit about Vision AutoML

## Step 3 - Demo Tracing

Visit the Stackdriver tracing UI and show it off for all its glory!

The app is actually sending traces directly via Spring Cloud GCP _and_ via ASM. Nifty eh?

## Step 4 - Demo Service Visibility

### Google Cloud

Visit the Anthos Service Mesh Services page in the cloud console to see the various graph, timeline and topology views.

Show SLO and SLI metrics at this stage

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

Go back to [crispcam.com](https://crispcam.com) and show it's broken!

### 6c - Fix it

Show the [yaml/destination-rule.yaml](yaml/destination-rule.yaml) file and show the circuit breaker configuration.

Now show [yaml/virtual-service-fixed.yaml](yaml/virtual-service-fixed.yaml) - point out how it now has a `timeout: 1s` parameter - if the request takes longer than that, intercept and return a gateway error.

Apply this change:

```
kubectl apply -f yaml/virtual-service-fixed.yaml
```

Now go back to [crispcam.com](https://crispcam.com) and observe the reviews are showing down, but the site is working properly.

## Step 7 - Cleanup

Apply the cleanup file so it works for everyone again:

```
kubectl apply -f yaml/virtual-service-lb.yaml
```
