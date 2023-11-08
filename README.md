log-pilot
=========

[![CircleCI](https://circleci.com/gh/shideqin/log-pilot.svg?style=svg)](https://circleci.com/gh/shideqin/log-pilot)
[![Go Report Card](https://goreportcard.com/badge/github.com/shideqin/log-pilot)](https://goreportcard.com/report/github.com/shideqin/log-pilot)

`log-pilot` is an awesome container log tool. With `log-pilot` you can collect logs from kubernetes hosts and send them to your centralized log system such as elasticsearch, graylog2, awsog and etc. `log-pilot` can collect not only kubernetes stdout but also log file that inside kubernetes containers.

Try it
======

Prerequisites:

- kubeernetes >= 1.20

```
# download log-pilot project
git clone git@github.com:shideqin/log-pilot.git
# build log-pilot image
cd log-pilot/ && ./build-image.sh
# quick start
cd quickstart/ && ./run
```

Then access kibana under the tips. You will find that tomcat's has been collected and sended to kibana.

Create index:
![kibana](quickstart/Kibana.png)

Query the logs:
![kibana](quickstart/Kibana2.png)

Quickstart
==========

### Run pilot

```
# log-pilot.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: log-pilot
  name: log-pilot
  namespace: default
spec:
  selector:
    matchLabels:
      k8s-app: log-pilot
  template:
    metadata:
      labels:
        k8s-app: log-pilot
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: node-role.kubernetes.io/master
                operator: DoesNotExist
              - key: type
                operator: NotIn
                values:
                - virtual-kubelet
      containers:
      - env:
        image: registry.cn-hangzhou.aliyuncs.com/ad-hub/log-pilot:v6.11.4.0
        imagePullPolicy: Always
        name: log-pilot
        resources: {}
        securityContext:
          capabilities:
            add:
            - SYS_ADMIN
        volumeMounts:
        - mountPath: /host
          name: root
          readOnly: true
        - mountPath: /etc/localtime
          name: localtime
        - mountPath: /etc/fluentd/fluentd.conf
          name: vol-uviis
          subPath: fluentd.conf
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: devops-admin
      serviceAccountName: devops-admin
      terminationGracePeriodSeconds: 30
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      volumes:
      - hostPath:
          path: /
          type: ""
        name: root
      - hostPath:
          path: /etc/localtime
          type: ""
        name: localtime
      - configMap:
          name: log-pilot-config
        name: vol-uviis

#kubectl
kubectl apply -f log-lipot.yaml
```

### Run applications whose logs need to be collected

Open a new terminal, run the application. With tomcat for example:

```
docker run -it --rm  -p 10080:8080 \
    -v /usr/local/tomcat/logs \
    --label aliyun.logs.catalina=stdout \
    --label aliyun.logs.access=/usr/local/tomcat/logs/localhost_access_log.*.txt \
    tomcat
```

Now watch the output of log-pilot. You will find that log-pilot get all tomcat's startup logs. If you access tomcat with your broswer, access logs in `/usr/local/tomcat/logs/localhost_access_log.\*.txt` will also be displayed in log-pilot's output.

More Info: [Fluentd Plugin](docs/fluentd/docs.md) and [Filebeat Plugin](docs/filebeat/docs.md)

Feature
========

- Support both [fluentd plugin](docs/fluentd/docs.md) and [filebeat plugin](docs/filebeat/docs.md). You don't need to create new fluentd or filebeat process for every kubernetes container.
- Support both stdout and log files. Either kubernetes log driver or logspout can only collect stdout.
- Declarative configuration. You need do nothing but declare the logs you want to collect.
- Support many log management: elastichsearch, graylog2, awslogs and more.
- Tags. You could add tags on the logs collected, and later filter by tags in log management.

Build log-pilot
===================

Prerequisites:

- Go >= 1.18

```
go get github.com/shideqin/log-pilot
cd $GOPATH/github.com/shideqin/log-pilot
# This will create a new docker image named log-pilot:latest
./build-image.sh fluentd
```

Contribute
==========

You are welcome to make new issues and pull reuqests.

