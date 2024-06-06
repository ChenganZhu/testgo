# Part 1
## Linter
The project uses golangci-lint. Install it locally and simply run
```
golangci-lint run
```
## Build and usage
```
go build
./myprogram -u https://www.google.fr -u https://news.ycombinator.com
```

# Part 2
## Docker
Build and run the container
```
docker build . -t myprogram
docker run -it myprogram -u "https://www.google.fr" -u "https://news.ycombinator.com"
```

Run Trivy security scan against the program
```
docker pull aquasec/trivy:0.52.0
docker run -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:0.52.0 image myprogram
```

As of today (Thursday, June 6, 2024), the scan finished with no known security vulnerability :
```
docker run -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:0.52.0 image myprogram
2024-06-06T15:01:34Z	INFO	Need to update DB
...
...
...
2024-06-06T15:01:37Z	INFO	Secret scanning is enabled
2024-06-06T15:01:37Z	INFO	If your scanning is slow, please try '--scanners vuln' to disable secret scanning
2024-06-06T15:01:37Z	INFO	Please see also https://aquasecurity.github.io/trivy/v0.52/docs/scanner/secret/#recommendation for faster secret detection
2024-06-06T15:01:37Z	INFO	Detected OS	family="debian" version="12.5"
2024-06-06T15:01:37Z	INFO	[debian] Detecting vulnerabilities...	os_version="12" pkg_num=3
2024-06-06T15:01:37Z	INFO	Number of language-specific files	num=1
2024-06-06T15:01:37Z	INFO	[gobinary] Detecting vulnerabilities...

myprogram (debian 12.5)
=======================
Total: 0 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 0, CRITICAL: 0)
```

## Kubernetes
Apply with kubectl the file `pod.yaml` for an example :
```
kubectl apply -f pod.yaml
```
A tail command has been added to the pod to make it sleep forever. To delete the pod use :
```
kubectl delete pod myprogram --force --grace-period 0
```
# Part 3
A simple CICD pipeline has been implemented with GitHubActions.
On any push to main branch, it does the following :
- lint the GO application
- build the GO application
- build and publish the container in a public repository to DockerHub

# Part 4
The file list_of_links contains :
```
http://tiktok.com
https://ads.faceBook.com.
https://sub.ads.faCebok.com
api.tiktok.com
Google.com.
aws.amazon.com
```

## First attempt to extract domains
```
cat list_of_links | tr '[A-Z]' '[a-z]' | sed 's/[.]$//' | sed 's/^http:[/]*//' | sed 's/^https:[/]*//' | awk -F"." '{print $(NF-1)"."$(NF)}' | sort | uniq
```

## Second attempt to extract domains

