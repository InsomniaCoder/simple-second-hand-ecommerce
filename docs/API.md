# API Documentation

## Base URL
```
http://localhost:8080
```

## Endpoints

### Health Check
Check if the API is running.

```http
GET /health
```

**Response**
```json
{
  "status": "ok"
}
```

---

### Create Item
Create a new item for sale.

```http
POST /items
Content-Type: application/json
```

**Request Body**
```json
{
  "title": "MacBook Pro 16-inch",
  "description": "2021 model, M1 Max, 32GB RAM, excellent condition",
  "price": 1999.99,
  "image": "https://example.com/macbook.jpg",
  "status": "available"
}
```

**Response** (201 Created)
```json
{
  "id": "3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c",
  "title": "MacBook Pro 16-inch",
  "description": "2021 model, M1 Max, 32GB RAM, excellent condition",
  "price": 1999.99,
  "image": "https://example.com/macbook.jpg",
  "status": "available",
  "created_at": "2024-01-05T14:30:00Z",
  "updated_at": "2024-01-05T14:30:00Z"
}
```

---

### List All Items
Get all items.

```http
GET /items
```

**Response** (200 OK)
```json
[
  {
    "id": "3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c",
    "title": "MacBook Pro 16-inch",
    "description": "2021 model",
    "price": 1999.99,
    "image": "https://example.com/macbook.jpg",
    "status": "available",
    "created_at": "2024-01-05T14:30:00Z",
    "updated_at": "2024-01-05T14:30:00Z"
  }
]
```

---

### Get Item
Get a specific item by ID.

```http
GET /items/{id}
```

**Response** (200 OK)
```json
{
  "id": "3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c",
  "title": "MacBook Pro 16-inch",
  "price": 1999.99,
  ...
}
```

**Error** (404 Not Found)
```json
{
  "error": "item not found"
}
```

---

### Update Item
Update an existing item (partial updates supported).

```http
PUT /items/{id}
Content-Type: application/json
```

**Request Body** (all fields optional)
```json
{
  "title": "MacBook Pro - Price Reduced",
  "price": 1799.99,
  "status": "booked"
}
```

**Response** (200 OK)
```json
{
  "id": "3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c",
  "title": "MacBook Pro - Price Reduced",
  "price": 1799.99,
  "status": "booked",
  "updated_at": "2024-01-05T15:00:00Z",
  ...
}
```

---

### Delete Item
Delete an item.

```http
DELETE /items/{id}
```

**Response** (200 OK)
```json
{
  "message": "item deleted"
}
```

---

## Status Values
- `available` - Item is available for purchase
- `booked` - Item is reserved/booked
- `sold` - Item has been sold

## Error Responses

### 400 Bad Request
```json
{
  "error": "invalid input: title is required and must not exceed 200 characters"
}
```

### 404 Not Found
```json
{
  "error": "item not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

## cURL Examples

### Create Item
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Desk Chair",
    "description": "Ergonomic office chair",
    "price": 150.00,
    "image": "https://example.com/chair.jpg"
  }'
```

### List Items
```bash
curl http://localhost:8080/items
```

### Get Item
```bash
curl http://localhost:8080/items/{id}
```

### Update Item
```bash
curl -X PUT http://localhost:8080/items/{id} \
  -H "Content-Type: application/json" \
  -d '{"price": 120.00, "status": "booked"}'
```

### Delete Item
```bash
curl -X DELETE http://localhost:8080/items/{id}
```
