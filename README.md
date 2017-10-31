# glasses

see all the domains that [traefik](https://traefik.io/) is handling.


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
# reading k8s config...
192.168.99.100 grafana.local # sad-chicken-grafana
192.168.99.100 prometheus.local # your-turkey-prometheus
```

write mode:
```
$ sudo -E glasses -write
# reading k8s config...
192.168.99.100 grafana.local # sad-chicken-grafana
192.168.99.100 prometheus.local # your-turkey-prometheus
```