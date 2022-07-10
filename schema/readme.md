# **Миграции для базы данных**

---
## ***Поднимаем базу в ```docker:```***
```
docker run --name=postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres
```

---

## ***Настраиваем ```миграции:```***

```
#Генерируем файлы и описываем в них талицы
migrate create -ext sql -dir . -seq init

#Создаем таблицы в базе данных
migrate -path . -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

#Откатываем изменения в базе данных
migrate -path . -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down
```

---