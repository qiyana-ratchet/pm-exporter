# PM Exporter
PM data exporter for Prometheus on Kubernetes Cluster

## Features

Kubernetes Cluster node에서 수신한 메트릭을 저장해 Prometheus의 http 요청에 따라 메트릭을 전송합니다.

#### Notice
1. 샘플 PM 데이터와 원본 데이터의 차이점
    - parse-me-*.xml의 맨 윗줄은 원본 파일에서 작업이 필요할 수 있습니다.  
    - 다만 일반적으로는 파싱할 때 무시하기 때문에 제대로 작동할 수 있습니다.  
1. 원본 데이터의 메트릭 ID 이름이 중복 사용되어 있는 사항  
    - EMS의 PM 데이터 형식을 바꿀수 있다면 간단하게 해결 될 수 있습니다.
    - 한 리스트 단위로 배열을 만들어 Validation하는 함수를 만들어 단체적으로 이름을 변경하는 함수를 만들 수 있습니다. 배열에서 중복된 이름이 있는지 검사해, 중복된 이름이 있다면 _1을 붙이거나 rename할 수 있습니다. 
    - Update 함수의 이중포문안에서 ch에 넣기 전에 함수를 호출하면 됩니다.
    - xml 원본 파일의 RACH.PreambleDed 데이터 부터는 이름이 잘 붙어서 나옵니다.
1. 클래스 접근
    - 중복되는 이름을 검사해서 바꿔주는 함수를 만들어서 사용
    - metricKey는 measDataFile.MeasData.MeasInfo[j].MeasType[i].Value 와 같이 읽을 수 있습니다.  

## Reference

- [Node Exporter](https://github.com/prometheus/node_exporter/blob/master/node_exporter.go) - 프로메테우스 노드 익스포터 깃허브 주소
- [client-golang/Prometheus](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#BuildFQName) - 프로메테우스 패키지 함수 설명
- [go-kit/Prometheus](https://pkg.go.dev/github.com/go-kit/kit/metrics/prometheus) - 프로메테우스 메트릭 타입 설명
- [Grafana](https://grafana.com/grafana/dashboards/14931) - Grafana 대시보드 

## Kubernetes
#### Installation

```sh
#GitLab에서 소스 가져오기
git clone http://218.233.172.197:90/thkim/pm-exporter
kubectl apply -f pm_exporter.yaml

kubectl exec -it testlog /bin/sh
pm-exporter
```

* 이미지 pull이 안될 경우 DockerHub에 이미지 올리기
```sh
docker build -t my_user_name/my_image:dev .
docker push my_user_name/my_image:dev 

kubectl apply -f pm_exporter.yaml
```

## Prometheus
#### Installation

Prometheus Helm chart의 value.yaml에서 Kubernetes Pod의 target ip를 변경합니다.  
helm install command
```sh
helm install prometheus my_prometheus/
```
![Screenshot_from_2021-08-26_15-14-06](/uploads/c09a5c92b0ed37fd4682ddd274133a1b/Screenshot_from_2021-08-26_15-14-06.png)

## Grafana
- Grafana를 이용해 시각화된 여러 메트릭 그래프를 동시에 모니터링 할 수 있습니다.  
- Grafana 대시보드는 레포지토리에 저장되어 있는 grafana-dashboard-1630025916657.json 파일을 import해서 사용할 수 있습니다.  
~~대시보드 Import ID는 14931 입니다~~


![Screenshot_from_2021-08-26_15-11-46](/uploads/c3ba66ec8cd47034187570e5bba04994/Screenshot_from_2021-08-26_15-11-46.png)
