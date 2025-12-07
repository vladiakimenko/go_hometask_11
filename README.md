# Запуск
```
go run cmd/server/main.go
```

# Примеры запросов curl

## Получить список задач
```
curl -X GET http://localhost:8080/tasks
```
**Статус код при успехе**: 200  

## Добавить новую задачу
```
curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Покормить кота"
         }'
```
**Обязательные поля**: `title`  
**Опциональные поля**: -   
**Статус код при успехе**: 201    

## Показать конкретную задачу
```
curl -X GET http://localhost:8080/tasks/1
```
*где '1' - это ID задачи*  
**Статус код при успехе**: 200    

## Зменить задачу
```
curl -X PUT http://localhost:8080/tasks/1 \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Помыть посуду",
           "done": true,
           "created_at": "2025-12-06T15:30:00Z"
         }'
```
*где '1' - это ID задачи*  
*непереданные поля выставляются в дефолт*  
**Обязательные поля**: `title`  
**Опциональные поля**: `done`, `created_at`   
**Статус код при успехе**: 200   

## Изменить задачу
```
curl -X PATCH http://localhost:8080/tasks/1 \
     -H "Content-Type: application/json" \
     -d '{
           "done": true
         }'
```
*где '1' - это ID задачи*  
*непереданные поля сохраняют своё значение*  
**Обязательные поля**: -  
**Опциональные поля**: `title`, `done`, `created_at`   
**Статус код при успехе**: 200   

## Удалить задачу
```
curl -X DELETE http://localhost:8080/tasks/1
```
*где '1' - это ID задачи*  
**Статус код при успехе**: 204