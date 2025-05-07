#Создание задачи (POST)
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
        "title": "Buy groceries",
        "description": "Milk, Bread, Eggs",
        "status": "new"
      }'
#Чтение списка задач (GET)
curl -i http://localhost:8080/tasks
#Чтение конкретной задачи (GET)
curl -i http://localhost:8080/tasks/1
#Обновление задачи (PUT)
curl -i -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
        "title": "Buy groceries and fruits",
        "description": "Milk, Bread, Eggs, Apples",
        "status": "in_progress"
      }'