[![CI/CD](https://github.com/tabularasa31/antibruteforce/actions/workflows/main.yml/badge.svg)](https://github.com/tabularasa31/antibruteforce/actions/workflows/main.yml)   [![Linters](https://github.com/tabularasa31/antibruteforce/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/tabularasa31/antibruteforce/actions/workflows/golangci-lint.yml)    [![Go Report Card](https://goreportcard.com/badge/github.com/tabularasa31/antibruteforce)](https://goreportcard.com/report/github.com/tabularasa31/antibruteforce)


# antibruteforce

## Общее описание
Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе.

Сервис вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.

Cервис используется только для server-server, т.е. скрыт от конечного пользователя.

## Алгоритм работы
Сервис ограничивает частоту попыток авторизации для различных комбинаций параметров, заданных в конфигурационном файле config.yml:
* не более N = 10 попыток в минуту для данного логина.
* не более M = 100 попыток в минуту для данного пароля (защита от обратного brute-force).
* не более K = 1000 попыток в минуту для данного IP (число большое, т.к. NAT).

Для подсчета и ограничения частоты запросов использован например алгоритм GCRA (aka leaky bucket)
https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm

Сервис поддерживает множество bucket-ов, по одному на каждый логин/пароль/ip.
Сами bucket-ы хранятся в Redis.

White/black листы содержат списки адресов сетей, которые обрабатываются более простым способом:
* Если входящий IP в whitelist, то сервис безусловно разрешает авторизацию (ok=true);
* Если в blacklist, то отклоняет (ok=false).
White/black листы хранятся в postgresql.

## Конфигурация
Конфигурация сервиса находится в файле config/config.yml .
Основные параметры конфигурации: N, M, K - лимиты по достижению которых, сервис считает попытку брутфорсом.

## Архитектура
Микросервис состоит из GRPC API, базы данных для хранения bucket'ов (Redis), black/white списков (Postgres) 
и command-line интерфейса взаимодействия с сервисом.

## Описание методов API

### Попытка авторизации
Запрос:
* login
* password
* ip

Ответ:
* ok (true/false) - сервис должен возвращать ok=true, если считает что запрос нормальный 
и ok=false, если считает что происходит bruteforce.

### Сброс bucket
Должен очистить bucket-ы соответствующие переданным login и ip.
* login
* ip

### Добавление IP в blacklist
* подсеть (IP + маска)

### Удаление IP из blacklist
* подсеть (IP + маска)

### Добавление IP в whitelist
* подсеть (IP + маска)

### Удаление IP из whitelist
* подсеть (IP + маска)

---

- Пример подсети: 192.1.1.0/25 - представляет собой адрес 192.1.1.0 с маской 255.255.255.128

---

## Command-Line интерфейс
Реализован command-line интерфейс для ручного администрирования сервиса.
Через CLI есть возможность вызвать запрос на проверку login, password и IP, 
сброс бакета, а так же управлять whitelist/blacklist-ом.

## Развертывание
Развертывание микросервиса должно осуществляется командой `make up`
в директории с проектом.



