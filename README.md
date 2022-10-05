# service_architecture_hw5

1. После запуска в Makefile надо получить адрес Grafana. Она находится в namespace prometheus-operator
kubectl describe -n prometheus-operator service - проверяем

**Вывод**:
`Name:              prometheus-operator-grafana
Namespace:         prometheus-operator
Labels:            app.kubernetes.io/instance=prometheus-operator
app.kubernetes.io/managed-by=Helm
app.kubernetes.io/name=grafana
app.kubernetes.io/version=9.0.5
helm.sh/chart=grafana-6.32.10
Annotations:       meta.helm.sh/release-name: prometheus-operator
meta.helm.sh/release-namespace: prometheus-operator
Selector:          app.kubernetes.io/instance=prometheus-operator,app.kubernetes.io/name=grafana
Type:              ClusterIP
IP Family Policy:  SingleStack
IP Families:       IPv4
IP:                10.109.222.21
IPs:               10.109.222.21
Port:              http-web  80/TCP
TargetPort:        3000/TCP
Endpoints:         172.17.0.5:3000
Session Affinity:  None
Events:            <none>`

# Пробрасываем порт во вне
kubectl port-forward -n prometheus-operator service/prometheus-operator-grafana 9000:80

# Заходим по адресу http://localhost:9000/ **Логин** admin **Пароль** prom-operator



# Problems
Не решена проблема с метриками ingress_nginx_controller_request
