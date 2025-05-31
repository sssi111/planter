#!/bin/bash

# Скрипт для добавления образцов растений в базу данных через API

# Проверка, запущен ли сервер
echo "Проверка доступности сервера..."
curl -s http://localhost:8080/plants > /dev/null
if [ $? -ne 0 ]; then
    echo "Ошибка: Сервер недоступен. Убедитесь, что сервер запущен на http://localhost:8080"
    exit 1
fi

echo "Сервер доступен. Начинаем добавление растений..."

# Путь к файлу с образцами растений
PLANTS_FILE="scripts/sample_plants.json"

# Проверка наличия файла
if [ ! -f "$PLANTS_FILE" ]; then
    echo "Ошибка: Файл $PLANTS_FILE не найден"
    exit 1
fi

# Используем Python для обработки JSON
# Проверяем, установлен ли Python
if ! command -v python3 &> /dev/null && ! command -v python &> /dev/null; then
    echo "Ошибка: Python не установлен. Пожалуйста, установите Python для запуска этого скрипта."
    exit 1
fi

# Определяем команду Python (python3 или python)
PYTHON_CMD="python3"
if ! command -v python3 &> /dev/null; then
    PYTHON_CMD="python"
fi

# Получаем количество растений и их данные с помощью Python
PLANTS_COUNT=$($PYTHON_CMD -c "import json; f=open('$PLANTS_FILE'); data=json.load(f); print(len(data)); f.close()")

echo "Найдено $PLANTS_COUNT растений для добавления"

# Обрабатываем каждое растение
for i in $(seq 0 $(($PLANTS_COUNT - 1))); do
    # Получаем имя растения
    PLANT_NAME=$($PYTHON_CMD -c "import json; f=open('$PLANTS_FILE'); data=json.load(f); print(data[$i]['name']); f.close()")
    
    echo "Добавление растения $PLANT_NAME..."
    
    # Получаем JSON для текущего растения
    PLANT=$($PYTHON_CMD -c "import json; f=open('$PLANTS_FILE'); data=json.load(f); print(json.dumps(data[$i])); f.close()")
    
    # Отправляем запрос
    echo "Отправляемый JSON:"
    echo "$PLANT" | $PYTHON_CMD -m json.tool
    
    RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$PLANT" \
        http://localhost:8080/admin/plants)
    
    echo "Ответ сервера:"
    echo "$RESPONSE"
    
    if [ $? -eq 0 ]; then
        # Проверяем ответ с помощью Python
        IS_SUCCESS=$($PYTHON_CMD -c "import json, sys;
try:
    data=json.loads('$RESPONSE');
    print('true' if 'id' in data else 'false')
except:
    print('false')")
        
        if [ "$IS_SUCCESS" = "true" ]; then
            PLANT_ID=$($PYTHON_CMD -c "import json; data=json.loads('$RESPONSE'); print(data['id'])")
            echo "Растение успешно добавлено с ID: $PLANT_ID"
        else
            echo "Ошибка при добавлении растения: $RESPONSE"
        fi
    else
        echo "Ошибка при выполнении запроса"
    fi
    
    # Небольшая пауза между запросами
    sleep 0.5
done

echo "Добавление растений завершено"