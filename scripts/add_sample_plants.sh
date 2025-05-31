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
    
    # Получаем JSON для текущего растения и сохраняем во временный файл
    TEMP_FILE=$(mktemp)
    $PYTHON_CMD -c "import json; f=open('$PLANTS_FILE'); data=json.load(f); f2=open('$TEMP_FILE', 'w'); json.dump(data[$i], f2, ensure_ascii=False); f.close(); f2.close()"
    
    # Отправляем запрос
    echo "Отправляемый JSON:"
    cat "$TEMP_FILE" | $PYTHON_CMD -m json.tool
    
    RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        --data-binary "@$TEMP_FILE" \
        http://localhost:8080/admin/plants)
    
    # Сохраняем код ответа
    STATUS_CODE=$?
    echo "Ответ сервера:"
    echo "$RESPONSE"
    
    # Проверяем ответ
    if [ $STATUS_CODE -eq 0 ]; then
        # Сохраняем ответ во временный файл для обработки
        RESPONSE_FILE=$(mktemp)
        echo "$RESPONSE" > "$RESPONSE_FILE"
        
        # Проверяем ответ с помощью Python
        RESULT=$($PYTHON_CMD -c "import json, sys;
try:
    with open('$RESPONSE_FILE') as f:
        data = json.load(f)
    if 'id' in data:
        print('SUCCESS:' + data['id'])
    else:
        print('ERROR:Неверный формат ответа сервера')
except Exception as e:
    print(f'ERROR:{str(e)}')")
        
        # Обрабатываем результат
        if [[ "$RESULT" == SUCCESS:* ]]; then
            PLANT_ID=${RESULT#SUCCESS:}
            echo "✅ Растение успешно добавлено с ID: $PLANT_ID"
        else
            ERROR_MSG=${RESULT#ERROR:}
            echo "❌ Ошибка при обработке ответа: $ERROR_MSG"
            echo "Полный ответ сервера:"
            echo "$RESPONSE"
        fi
        
        rm "$RESPONSE_FILE"
    else
        echo "❌ Ошибка при выполнении запроса (код $STATUS_CODE)"
    fi
    
    rm "$TEMP_FILE"
    
    # Пауза между запросами и разделитель для читаемости
    echo ""
    echo "--------------------------------------------------"
    sleep 1
done

echo "Добавление растений завершено"