Lamoda Test Task
===========
____

- Перед запуском контейнеров нужно заполнить .env
- ### Документация по запросам представлена ниже в формате CURL запросов. 
- Если требуется развернуть в локальной базе данных, нужно создать таблицу и импортировать туда db.sql

p.s - При разработке использовал постман,но коллекция без документации. Возможно будет удобнее использовать её для проверки)

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/18001533-99cabde6-8dd4-469f-ab28-1c6ef6074068?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D18001533-99cabde6-8dd4-469f-ab28-1c6ef6074068%26entityType%3Dcollection%26workspaceId%3Deec3497c-8cb7-4583-a9d9-2afefd140d53)

#### Первый запуск

1. Скопировать `.env.example` в файл `.env`. При желании изменить в нём значения.
2. Запустить команду `make run`

----
#### Запуск тестов
1. Запустить команду `make test`
----
#### Посмотреть покрытие тестами
1. Запустить команду `make coverage`
----
#### Запуск, перезапуск и остановка докер контейнеров
- Для запуска, команда `make up`
- Для остановки, команда `make down`
- Для перезапуска, команда `make restart`

-----
### Curl команды и результат

Мы всегда получаем json с полями 
1. Обязательное поле `code` с http кодом результата
2. Поле `message`. Если результат не подразумевает возврата полезной нагрузки, или произошла ошибка
3. Поле `data`. Если возвращается полезная нагрузка.
---
##### good/all 
Команда `curl --location '127.0.0.1:8080/goods/all'`

Возвращает все товары из таблицы в массиве:

Результат

    {
        "code": 200,
        "data": [
            {
                "id": 6,
                "name": "TestAddedFromAPI",
                "size": "XS",
                "uniq_code": 565
            }
        ]
    }
---
##### goods/remains
Команда `curl --location '127.0.0.1:8080/goods/remains'`

Возвращает объект с ключом `uniq_code` содержащий в себе:
1. `name` - имя товара
2. `size` - размер
3. `storage_avalable` - объект где каждому id склада соответствует количество товара с этим `uniq_code` доступное на этом складе

Результат

    {
        "code": 200,
        "data": {
            "1": {
                "name": "Test1",
                "size": "L",
                "storage_available": {
                    "1": 15,
                    "3": 10
                }
            }
        }
    }
---
##### goods/reserve
Команда 

    curl --location 'http://127.0.0.1:8080/goods/reserve' \
        --header 'Content-Type: application/json' \
        --data '[
            {
                "uniq_code":1,
                "count":200
            },
            {
                "uniq_code":2,
                "count":100
            }
        ]'
Входные значения:
1. `uniq_code` - уникальный код
2. `count` - сколько требуется зарезервировать товара

Возвращает массив объектов содержащих в себе:
1. `uniq_code` - уникальный код товара
2. `storage_avalable` - массив объектов где указано сколько этого товара было зарезервировано на конкретном складе
3. `additional_info` - Сопровождающая информация, если требуется. (В примере не было доступных товаров на всех доступных складах)

Результат

    {
        "code": 200,
        "data": [
            {
                "uniq_code": 1,
                "storages": [
                    {
                        "reserved": 15,
                        "storage": 1
                    },
                    {
                        "reserved": 5,
                        "storage": 3
                    }
                ]
            },
            {
                "uniq_code": 2,
                "storages": [],
                "additional_info": "Can't reserve this good"
            }
        ]
    }

---
##### goods/release
Команда

    curl --location 'http://127.0.0.1:8080/goods/release' \
        --header 'Content-Type: application/json' \
        --data '[
            {
                "uniq_code":1,
                "count":1
            },
            {
                "uniq_code":2,
                "count":1
            }
        ]'
Входные значения:
1. `uniq_code` - уникальный код товара
2. `count` - сколько требуется освободить товара

Возвращает массив объектов содержащих в себе:
1. `uniq_code` - уникальный код товара
2. `additional_info` - Сопровождающая информация. OK - успешно освобождён товар из резерва, иначе ошибка.

Результат

    {
        "code": 200,
        "data": [
            {
                "uniq_code": 1,
                "additional_info": "OK"
            },
            {
                "uniq_code": 2,
                "additional_info": "can't release this good"
            }
        ]
    }

---
##### goods/add
Команда

    curl --location --request PUT '127.0.0.1:8080/goods/add' \
        --header 'Content-Type: application/json' \
        --data '{
            "name":"TestAddedFromAPI",
            "size":"XS",
            "uniq_code":565
        }'
Входные значения:
1. `name` - название товара
2. `size` - размер товара
3. `uniq_code` - уникальный код товара

Возвращает id добавленной записи

Результат

    {
        "code": 200,
        "data": 10
    }
----
##### goods/delete
Команда

    curl --location --request DELETE '127.0.0.1:8080/goods/delete' \
        --header 'Content-Type: application/json' \
        --data '{
            "uniq_code":565
        }'
Входные значения:
1. `uniq_code` - уникальный код товара

Возвращает сообщение OK, если успешно.

Результат

    {
        "code": 200,
        "message": "OK"
    }
----
##### storages/all
Команда `curl --location 'http://127.0.0.1:8080/storages/all'`

Возвращает все склады из таблицы в массиве

Результат

    {
        "code": 200,
        "data": [
            {
                "id": 1,
                "name": "TestStore",
                "available": true
            },
            {
                "id": 2,
                "name": "Storage2",
                "available": false
            }
        ]
    }
---
##### storages/available
Команда `curl --location 'http://127.0.0.1:8080/storages/available'`

Возвращает все ДОСТУПНЫЕ склады из таблицы в массиве

Результат

    {
        "code": 200,
        "data": [
            {
                "id": 1,
                "name": "TestStore",
                "available": true
            }
        ]
    }
---
##### storages/add
Команда

    curl --location --request PUT '127.0.0.1:8080/storages/add' \
        --header 'Content-Type: application/json' \
        --data '{
            "name":"TestAddedFromAPI",
            "available": false
        }'

Входные значения:
1. `name` - название склада
2. `available` - доступность склада

Возвращает id добавленной записи

Результат

    {
        "code": 200,
        "data": 10
    }
----
##### storages/delete
Команда

    curl --location --request DELETE '127.0.0.1:8080/storages/delete' \
        --header 'Content-Type: application/json' \
        --data '{
            "id":6
        }'
Входные значения:
1. `id` - id склада

Возвращает сообщение OK, если успешно.

Результат

    {
        "code": 200,
        "message": "OK"
    }
----
##### storages/access
Команда

    curl --location '127.0.0.1:8080/storages/access' \
        --header 'Content-Type: application/json' \
        --data '{
            "id":6,
            "available":true
        }'
Входные значения:
1. `id` - id склада
2. `available` - новое состояние склада

Возвращает сообщение OK, если успешно.

Результат

    {
        "code": 200,
        "message": "OK"
    }

