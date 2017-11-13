# glasses

see all the [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) resources from [k8s](https://kubernetes.io/).


## Installation 
```
go get -u github.com/wakeful/glasses
```

## Usage

help mode:
```
$ glasses -h
Usage of glasses:
  -host-file string
        host file location (default "/etc/hosts")
  -write
        rewrite host file?
```

dry-run mode:
```
$ glasses
# reading k8s ingress resource...
192.168.99.100 grafana.local # sad-chicken-grafana
192.168.99.100 prometheus.local # your-turkey-prometheus
```

populate your `/etc/hosts` file:
```
$ sudo -E glasses -write
# reading k8s ingress resource...
192.168.99.100 grafana.local # sad-chicken-grafana
192.168.99.100 prometheus.local # your-turkey-prometheus

$ cat /etc/hosts
# generated using glasses start #
192.168.99.100 grafana.local # sad-chicken-grafana
192.168.99.100 prometheus.local # your-turkey-prometheus

# generated using glasses end #
```