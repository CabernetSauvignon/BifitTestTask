# Тестовое задание на позицию Backend developer (Go) в компанию БИФИТ

Был реализован WEB-сервис, который принимает и отправляет данные по HTTP протоколу в формате JSON. 

## Поддерживаемые запросы
1. Запрос возвращает список всех задач 
    'GET' 'http://localhost:8080/jobs \'
    -H 'accept: application/json' \
2. Запрос создаёт новую задачу с заданным именем
    'POST' 'http://localhost:8080/jobs \'
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{"name": "jobname"}'
3. Запрос изменяет задачу с заданным именем, добавляя ей дополнительные 5 минут до перехода задачи в статус "done"
    'PUT' 'http://localhost:8080/jobs'
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{"name": "jobname"}'
