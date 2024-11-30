# Описание сервисов

## - API Gateway - точка входа и агрегатор.
### Создание заказа
```
POST /orders
Headers:
Authorization: Bearer <JWT>
Body:
{
    "items": [
        {"id": "item_id", "quantity": 1}
    ],
    "promo_code": "PROMO123",
    "delivery": {
        "type": "pickup|courier",
        "address": "optional_for_courier",
        "store_id": "optional_for_pickup"
    },
    "city_id": "city_uuid"
}
Response:
{
    "order_id": "uuid",
    "status": "created"
}
```
### Обновление способа оплаты заказа
```
PATCH /orders/{order_id}/payment
Body:
{
    "type": "cash|online|card_on_delivery|cash_on_delivery",
}
Response:
{
    "status": "updated"
}
```
## - User - хранит информацию о клиентах.
### Получение информации о пользователе:
```
POST /auth
 Body:
{
"token": "jwt-token"
}
Response:
{
"user_id": "uuid",
}
```

## - Order - сервис заказов. Отвечает за агрегацию и валидацию данных заказа, обновление статусов заказа.

### Создание заказа:
```
POST /orders
Body:
{
    "items": [
        {"id": "item_id", "quantity": 1}
    ],
    "promo_code": "PROMO123",
    "delivery": {
        "type": "pickup|courier",
        "address": "optional_for_courier",
        "store_id": "optional_for_pickup"
    },
    "city_id": "city_uuid",
    "user_id": "uuid"
}
Response:
{
    "order_id": "uuid",
    "status": "created"
}
```
### Обновление способа оплаты заказа:
```
PATCH /orders/{order_id}/payment
Body:
{
    "type": "cash|online|card_on_delivery|cash_on_delivery",
}
Response:
{
    "status": "updated"
}
```
### Получение информации о заказе:
```
GET /orders/{order_id}
Response:
{
    "order_id": "uuid",
    "status": "not_paid|paid|being_delivered|delivered",
    "payment_status": "not_paid|paid",
    "items": [
        {"id": "item_id", "quantity": 1, "price": 100}
    ],
    "delivery": {
        "type": "pickup|courier",
        "address": "optional_for_courier",
        "store_id": "optional_for_pickup"
    }
}
```

## - Store - сервис-агрегатор, инкапсулирующий логику меню, стоимости товаров и промо.

```
### Подтверждение заказа:

POST /validation
Body:
{
    "city_id": "city_uuid",
    "items": [
        {"id": "item_id", "quantity": 1}
    ],
    "promo_code": "PROMO123"
}
Response:
{
    "is_valid": true,
    "total_price": 1000,
    "estimated_ready_time": "2024-11-26T12:00:00Z"
}
```

## - Delivery - сервис курьерской доставки. Отвечает за координацию курьеров и верификацию доступности доставки.

```
### Проверка доступности доставки:

POST /availability
Body:
{
    "address": "delivery address",
    "city_id": "city_uuid"
}
Response:
{
    "is_available": true,
    "estimated_time": "2024-11-26T14:00:00Z"
}
```

## - Payment - сервис оплат. Инкапсулирует логику оплаты и походов во внешние сервисы оплат.

```
### Создание платежа:

POST /payments
Body:
{
    "order_id": "uuid",
    "amount": 1000,
    "payment_method": "online|cash_on_delivery|card_on_delivery"
}
Response:
{
    "payment_id": "uuid",
    "status": "success|failure"
}

### Проверка статуса платежа:

GET /payments/{payment_id}
Response:
{
    "payment_id": "uuid",
    "status": "success|failure|processing"
}
```

# Пользовательский сценарий оформления заказа

### 1. API Gateway отправляет запрос с JWT из запроса в сервис User

### 2. API Gateway формирует запрос в сервис Order, состоящий из:

### - Заказ (позиции меню, их количество, промокод)

### - Покупатель (уникальный идентификатор)

### - Доставка (способ доставки, адрес - для курьерской доставки, уникальный идентификатор магазина - для самовывоза)

### - Уникальный идентификатор населенного пункта пользователя

### 3. Сервис Order принимает запрос. Делает запрос в сервис Store за подтверждением состава заказа.

### 4. Сервис Store подтвержает, что состав и стоимость заказа в выбранном населенном пункте равен заказу из запроса. Если передан промокод, проверяет его наличие в базе. Если передан флаг самовывоза, то считает примерную дату готовности заказа. Возвращает флаг в сервис Order, если заказ провалидирован.

### 5. Сервис Order принимает ответ валидации заказа. Если ответ положительный, то, в зависимости от выбранного типа доставки, либо сразу сохраняем заказ в БД для самовывоза со статусом “not_paid” и датой готовности, либо, если выбрана курьерская доставка, верифицирует возможность курьерской доставки из сервиса Delivery.

### 6. После создания заказа, Order возвращает 201 ОК и номер заказа. На клиентской стороне даем возможность оплатить заказ.

### 7. Клиент отправляет запрос с номером заказа и способом оплаты в сервис Order.

### 8. Обновляем статусы оплаты и доставки

### - Если способ оплаты - “cash” и способ доставки - “pickup”, то подразумеваем, что клиент оплачивает заказ в кафе. Обновляем статус оплаты заказа в “not_paid”.

### - Если способ оплаты - “online”, то отправляем запрос на оплату в сервис Payment. В случае успеха, меняем статус оплаты заказа на “paid” и публикуем событие в топик “order_to_be_delivered”.

### - Если способ оплаты “card_on_delivery” или “cash_on_delivery”, то обновляем статус оплаты заказа в “not_paid” и публикуем событие в топик “order_to_be_delivered”.

### 9. Сервис Delivery читает события топика “order_to_be_delivered” и начинает процесс доставки. Публикует событие в топик “order_being_delivered” для обновления статуса заказа сервисом Order на “being_delivered”.

### 10. При доставке, если статус оплаты заказа “not_paid”, курьер принимает оплату и через приложение проставляет заказу статус оплаты “paid” и завершает заказ.

### 11. После доставки сервис Delivery публикует событие в топик “order_delivered”.

### 12. Сервис Order меняет статус заказа на “delivered”.

### 13. Конечное состояние заказа в БД - статус оплаты “paid”, статус доставки “delivered”.
