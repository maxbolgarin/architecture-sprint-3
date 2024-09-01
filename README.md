# Разбор кейса

Компания называется "Smart Home"


## Текущее решение

### Описание

Smart Home — приложение, позволяющее пользователям управлять датчиками умного дома и отслеживать их состояние.

В данный момент приложение представляет из себя монолит, написанный на Java с применением фреймворка Spring. В качестве базы данных используется PostgreSQL.

Сборка выполняется с помощью Maven, приложение помещается в Docker-image. В GitHub настроен пайплайн, который билдит образ и отправляет в registry. Раскатка выполняется в Kubernetes, для этого написан helm chart и terraform.

### Функциональность

Монолит предоставляет 6 эндпоинтов для контроля сенсоров и получения информации:

* GET `/api/heating/{id}` — получение информации о системе отопления по ID;
* PUT `/api/heating/{id}` —  обновление информации о системе отопления по ID;
* POST `/api/heating/{id}/turn-on` — включение системы отопления по ID;
* POST `/api/heating/{id}/turn-off` — отключене системы отопления по ID;
* POST `/api/heating/{id}/set-temperature` - установка целевой температуры для системы отопления по ID;
* GET `/api/heating/{id}/current-temperature` - получение текущей температуры системы отопления по ID.

В базе данных есть две таблицы - `heating_systems`, которая используется в описанных выше эндпоинтах, и `temperature_sensors`, которая не используется в данный момент.

### Анализ функциональности

В системе в данный момент есть один домен — heating_systems, нагревательные системы. Датчики температур хоть и заявлены в описании, но по факту не реализованы.

Плюсы реализации: написан рабочий код для минимального продукта, несящего ценность. Код небольшой, всего один сервис, его легко понять и разворачивать.

Минусы реализации: не работает заявленный функционал, отсутствует возможность масштабирования на требование добавления датчиков с "будущее неуточненное поведение", так как текущая архитектура "хардкодит" тип датчика.


## Целевое решение


# Базовая настройка

## Запуск minikube

[Инструкция по установке](https://minikube.sigs.k8s.io/docs/start/)

```bash
minikube start
```


## Добавление токена авторизации GitHub

[Получение токена](https://github.com/settings/tokens/new)

```bash
kubectl create secret docker-registry ghcr --docker-server=https://ghcr.io --docker-username=<github_username> --docker-password=<github_token> -n default
```


## Установка API GW kusk

[Install Kusk CLI](https://docs.kusk.io/getting-started/install-kusk-cli)

```bash
kusk cluster install
```


## Настройка terraform

[Установите Terraform](https://yandex.cloud/ru/docs/tutorials/infrastructure-management/terraform-quickstart#install-terraform)


Создайте файл ~/.terraformrc

```hcl
provider_installation {
  network_mirror {
    url = "https://terraform-mirror.yandexcloud.net/"
    include = ["registry.terraform.io/*/*"]
  }
  direct {
    exclude = ["registry.terraform.io/*/*"]
  }
}
```

## Применяем terraform конфигурацию 

```bash
cd terraform
terraform apply
```

## Настройка API GW

```bash
kusk deploy -i api.yaml
```

## Проверяем работоспособность

```bash
kubectl port-forward svc/kusk-gateway-envoy-fleet -n kusk-system 8080:80
curl localhost:8080/hello
```


## Delete minikube

```bash
minikube delete
```
