### Запуск

```shell
kubectl create namespace saga
helm install saga . -n saga

POST/GET http://arch.homework/orders
```

### Взаимодействие с приложением осуществляется через сервис Order.

## В качестве оркестратора выступает сервис Orchestrator,
## который выполняет в заданном порядке шаги саги в сервисы billing, warehouse и delivery