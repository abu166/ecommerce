# API Testing Documentation

## 1. User Endpoints

### 1.1 POST /users/register - Register a New User (Success)

**Description:** Register a user with valid details.
- **Method:** POST
- **URL:** `{{base_url}}/users/register`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "username": "alice",
  "password": "secret123",
  "email": "alice@example.com"
}
```
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "<uuid>",
  "username": "alice",
  "email": "alice@example.com"
}
```
- **Postman Tests (JavaScript):**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains id, username, email", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.be.a("string");
    pm.expect(jsonData.username).to.equal("alice");
    pm.expect(jsonData.email).to.equal("alice@example.com");
});
```
- **Notes:** Save the id from the response for login tests.

### 1.2 POST /users/register - Duplicate Username (Failure)

**Description:** Attempt to register a user with an existing username.
- **Method:** POST
- **URL:** `{{base_url}}/users/register`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "username": "alice",
  "password": "anotherpass",
  "email": "alice2@example.com"
}
```
- **Expected Response:**
    - **Status:** 500 Internal Server Error
    - **Body:**
```json
{
  "error": "failed to register user: pq: duplicate key value violates unique constraint \"users_username_key\""
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 500", function () {
    pm.response.to.have.status(500);
});
pm.test("Response contains error message", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.error).to.include("duplicate key");
});
```
- **Notes:** Run after the first registration to test the unique username constraint.

### 1.3 POST /users/register - Invalid Input (Failure)

**Description:** Register with missing fields.
- **Method:** POST
- **URL:** `{{base_url}}/users/register`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "username": "",
  "password": "",
  "email": ""
}
```
- **Expected Response:**
    - **Status:** 500 Internal Server Error
    - **Body:**
```json
{
  "error": "failed to register user: some database error"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 500", function () {
    pm.response.to.have.status(500);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.be.a("string");
});
```
- **Notes:** Tests input validation (though the app could improve validation).

### 1.4 POST /users/login - Successful Login

**Description:** Authenticate a registered user.
- **Method:** POST
- **URL:** `{{base_url}}/users/login`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "username": "alice",
  "password": "secret123"
}
```
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "token": "<uuid>"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains token", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.token).to.be.a("string");
    pm.environment.set("token", jsonData.token); // Store token
});
```
- **Notes:** Save the token in the token environment variable for authenticated requests. Run after registering "alice".

### 1.5 POST /users/login - Invalid Credentials (Failure)

**Description:** Attempt login with wrong password.
- **Method:** POST
- **URL:** `{{base_url}}/users/login`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "username": "alice",
  "password": "wrongpass"
}
```
- **Expected Response:**
    - **Status:** 401 Unauthorized
    - **Body:**
```json
{
  "error": "authentication failed: invalid credentials"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 401", function () {
    pm.response.to.have.status(401);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.equal("authentication failed: invalid credentials");
});
```

### 1.6 POST /users/login - Non-existent User (Failure)

**Description:** Attempt login with a non-existent username.
- **Method:** POST
- **URL:** `{{base_url}}/users/login`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "username": "bob",
  "password": "secret123"
}
```
- **Expected Response:**
    - **Status:** 401 Unauthorized
    - **Body:**
```json
{
  "error": "authentication failed: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 401", function () {
    pm.response.to.have.status(401);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("record not found");
});
```

### 1.7 GET /users/:id - Get User Profile (Success)

**Description:** Retrieve the profile of a registered user.
- **Method:** GET
- **URL:** `{{base_url}}/users/{{token}}`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "{{token}}",
  "username": "alice",
  "email": "alice@example.com"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains correct user details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.equal(pm.environment.get("token"));
    pm.expect(jsonData.username).to.equal("alice");
    pm.expect(jsonData.email).to.equal("alice@example.com");
});
```
- **Notes:** Use the token from the login response.

### 1.8 GET /users/:id - Invalid Token (Failure)

**Description:** Attempt to get a profile with an invalid token.
- **Method:** GET
- **URL:** `{{base_url}}/users/invalid-token`
- **Headers:**
    - Authorization: invalid-token
- **Body:** None
- **Expected Response:**
    - **Status:** 401 Unauthorized
    - **Body:**
```json
{
  "error": "invalid token"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 401", function () {
    pm.response.to.have.status(401);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.equal("invalid token");
});
```

### 1.9 GET /users/:id - Missing Token (Failure)

**Description:** Attempt to get a profile without an Authorization header.
- **Method:** GET
- **URL:** `{{base_url}}/users/some-id`
- **Headers:** None
- **Body:** None
- **Expected Response:**
    - **Status:** 401 Unauthorized
    - **Body:**
```json
{
  "error": "missing token"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 401", function () {
    pm.response.to.have.status(401);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.equal("missing token");
});
```

## 2. Product Endpoints

### 2.1 POST /products - Create Product (Success)

**Description:** Create a new product with valid details.
- **Method:** POST
- **URL:** `{{base_url}}/products`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "name": "Tablet",
  "category": "Electronics",
  "stock": 10,
  "price": 499.99
}
```
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "<uuid>",
  "name": "Tablet",
  "category": "Electronics",
  "stock": 10,
  "price": 499.99
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains product details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.be.a("string");
    pm.expect(jsonData.name).to.equal("Tablet");
    pm.expect(jsonData.category).to.equal("Electronics");
    pm.expect(jsonData.stock).to.equal(10);
    pm.expect(jsonData.price).to.equal(499.99);
    pm.environment.set("product_id", jsonData.id); // Store for later tests
});
```
- **Notes:** Save the id as product_id for order and product tests.

### 2.2 POST /products - Invalid Input (Failure)

**Description:** Attempt to create a product with negative stock.
- **Method:** POST
- **URL:** `{{base_url}}/products`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "name": "Invalid Product",
  "category": "Electronics",
  "stock": -5,
  "price": 100.00
}
```
- **Expected Response:**
    - **Status:** 200 OK (Note: The app accepts negative stock, which might be a bug to fix later)
    - **Body:**
```json
{
  "id": "<uuid>",
  "name": "Invalid Product",
  "category": "Electronics",
  "stock": -5,
  "price": 100.00
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains product details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.stock).to.equal(-5);
});
```
- **Notes:** Ideally, the app should validate stock ≥ 0 (consider adding validation).

### 2.3 POST /products - Missing Token (Failure)

**Description:** Attempt to create a product without authentication.
- **Method:** POST
- **URL:** `{{base_url}}/products`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "name": "Tablet",
  "category": "Electronics",
  "stock": 10,
  "price": 499.99
}
```
- **Expected Response:**
    - **Status:** 401 Unauthorized
    - **Body:**
```json
{
  "error": "missing token"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 401", function () {
    pm.response.to.have.status(401);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.equal("missing token");
});
```

### 2.4 GET /products/:id - Get Product (Success)

**Description:** Retrieve a product by ID.
- **Method:** GET
- **URL:** `{{base_url}}/products/{{product_id}}`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "{{product_id}}",
  "name": "Tablet",
  "category": "Electronics",
  "stock": 10,
  "price": 499.99
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains correct product details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.equal(pm.environment.get("product_id"));
    pm.expect(jsonData.name).to.equal("Tablet");
    pm.expect(jsonData.category).to.equal("Electronics");
    pm.expect(jsonData.stock).to.equal(10);
    pm.expect(jsonData.price).to.equal(499.99);
});
```
- **Notes:** Use product_id from the create product test.

### 2.5 GET /products/:id - Non-existent Product (Failure)

**Description:** Attempt to get a non-existent product.
- **Method:** GET
- **URL:** `{{base_url}}/products/non-existent-id`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 404 Not Found
    - **Body:**
```json
{
  "error": "product not found: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 404", function () {
    pm.response.to.have.status(404);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("product not found");
});
```

### 2.6 PATCH /products/:id - Update Product (Success)

**Description:** Update an existing product's details.
- **Method:** PATCH
- **URL:** `{{base_url}}/products/{{product_id}}`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "name": "Updated Tablet",
  "category": "Electronics",
  "stock": 20,
  "price": 599.99
}
```
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "{{product_id}}",
  "name": "Updated Tablet",
  "category": "Electronics",
  "stock": 20,
  "price": 599.99
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains updated product details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.equal(pm.environment.get("product_id"));
    pm.expect(jsonData.name).to.equal("Updated Tablet");
    pm.expect(jsonData.stock).to.equal(20);
    pm.expect(jsonData.price).to.equal(599.99);
});
```

### 2.7 PATCH /products/:id - Non-existent Product (Failure)

**Description:** Attempt to update a non-existent product.
- **Method:** PATCH
- **URL:** `{{base_url}}/products/non-existent-id`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "name": "Tablet",
  "category": "Electronics",
  "stock": 10,
  "price": 499.99
}
```
- **Expected Response:**
    - **Status:** 500 Internal Server Error
    - **Body:**
```json
{
  "error": "failed to update product: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 500", function () {
    pm.response.to.have.status(500);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("record not found");
});
```

### 2.8 DELETE /products/:id - Delete Product (Success)

**Description:** Delete an existing product.
- **Method:** DELETE
- **URL:** `{{base_url}}/products/{{product_id}}`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "message": "product deleted"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response confirms deletion", function () {
    pm.expect(pm.response.json().message).to.equal("product deleted");
});
```
- **Notes:** Create a new product before testing deletion.

### 2.9 DELETE /products/:id - Non-existent Product (Failure)

**Description:** Attempt to delete a non-existent product.
- **Method:** DELETE
- **URL:** `{{base_url}}/products/non-existent-id`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 500 Internal Server Error
    - **Body:**
```json
{
  "error": "failed to delete product: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 500", function () {
    pm.response.to.have.status(500);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("record not found");
});
```

### 2.10 GET /products - List Products (Success)

**Description:** List products with pagination and category filter.
- **Method:** GET
- **URL:** `{{base_url}}/products?page=1&page_size=10&category=Electronics`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "products": [
    {
      "id": "<uuid>",
      "name": "Tablet",
      "category": "Electronics",
      "stock": 10,
      "price": 499.99
    }
  ],
  "total": 1
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains products array and total", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.products).to.be.an("array");
    pm.expect(jsonData.total).to.be.a("number");
    if (jsonData.products.length > 0) {
        pm.expect(jsonData.products[0].category).to.equal("Electronics");
    }
});
```
- **Notes:** Create at least one product in "Electronics" category first.

### 2.11 GET /products - Empty List (Success)

**Description:** List products when none exist.
- **Method:** GET
- **URL:** `{{base_url}}/products?page=1&page_size=10`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "products": [],
  "total": 0
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains empty products array", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.products).to.be.an("array").that.is.empty;
    pm.expect(jsonData.total).to.equal(0);
});
```
- **Notes:** Clear the database or ensure no products exist.

## 3. Order Endpoints

### 3.1 POST /orders - Create Order (Success)

**Description:** Create an order with valid product IDs.
- **Method:** POST
- **URL:** `{{base_url}}/orders`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 2
    }
  ]
}
```
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "<uuid>",
  "user_id": "{{token}}",
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 2
    }
  ],
  "status": "pending",
  "total": 999.98
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains order details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.be.a("string");
    pm.expect(jsonData.user_id).to.equal(pm.environment.get("token"));
    pm.expect(jsonData.items).to.be.an("array").with.lengthOf(1);
    pm.expect(jsonData.items[0].product_id).to.equal(pm.environment.get("product_id"));
    pm.expect(jsonData.items[0].quantity).to.equal(2);
    pm.expect(jsonData.status).to.equal("pending");
    pm.expect(jsonData.total).to.equal(999.98);
    pm.environment.set("order_id", jsonData.id); // Store for later tests
});
```
- **Notes:** Ensure a product exists with sufficient stock (e.g., 10 units).

### 3.2 POST /orders - Insufficient Stock (Failure)

**Description:** Attempt to create an order with excessive quantity.
- **Method:** POST
- **URL:** `{{base_url}}/orders`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 100
    }
  ]
}
```
- **Expected Response:**
    - **Status:** 500 Internal Server Error
    - **Body:**
```json
{
  "error": "failed to create order: insufficient stock"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 500", function () {
    pm.response.to.have.status(500);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.equal("failed to create order: insufficient stock");
});
```

### 3.3 POST /orders - Non-existent Product (Failure)

**Description:** Attempt to create an order with an invalid product ID.
- **Method:** POST
- **URL:** `{{base_url}}/orders`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "items": [
    {
      "product_id": "non-existent-id",
      "quantity": 1
    }
  ]
}
```
- **Expected Response:**
    - **Status:** 500 Internal Server Error
    - **Body:**
```json
{
  "error": "failed to create order: product not found: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 500", function () {
    pm.response.to.have.status(500);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("product not found");
});
```

### 3.4 POST /orders - Missing Token (Failure)

**Description:** Attempt to create an order without authentication.
- **Method:** POST
- **URL:** `{{base_url}}/orders`
- **Headers:**
    - Content-Type: application/json
- **Body (raw, JSON):**
```json
{
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 1
    }
  ]
}
```
- **Expected Response:**
    - **Status:** 401 Unauthorized
    - **Body:**
```json
{
  "error": "missing token"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 401", function () {
    pm.response.to.have.status(401);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.equal("missing token");
});
```

### 3.5 GET /orders/:id - Get Order (Success)

**Description:** Retrieve an order by ID.
- **Method:** GET
- **URL:** `{{base_url}}/orders/{{order_id}}`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "{{order_id}}",
  "user_id": "{{token}}",
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 2
    }
  ],
  "status": "pending",
  "total": 999.98
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains correct order details", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.equal(pm.environment.get("order_id"));
    pm.expect(jsonData.user_id).to.equal(pm.environment.get("token"));
    pm.expect(jsonData.items[0].product_id).to.equal(pm.environment.get("product_id"));
    pm.expect(jsonData.status).to.equal("pending");
    pm.expect(jsonData.total).to.equal(999.98);
});
```

### 3.6 GET /orders/:id - Non-existent Order (Failure)

**Description:** Attempt to get a non-existent order.
- **Method:** GET
- **URL:** `{{base_url}}/orders/non-existent-id`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 404 Not Found
    - **Body:**
```json
{
  "error": "order not found: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 404", function () {
    pm.response.to.have.status(404);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("order not found");
});
```

### 3.7 PATCH /orders/:id - Update Order Status (Success)

**Description:** Update an order's status to "shipped".
- **Method:** PATCH
- **URL:** `{{base_url}}/orders/{{order_id}}`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "status": "shipped"
}
```
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "id": "{{order_id}}",
  "user_id": "{{token}}",
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 2
    }
  ],
  "status": "shipped",
  "total": 999.98
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains updated status", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.id).to.equal(pm.environment.get("order_id"));
    pm.expect(jsonData.status).to.equal("shipped");
});
```

### 3.8 PATCH /orders/:id - Non-existent Order (Failure)

**Description:** Attempt to update a non-existent order.
- **Method:** PATCH
- **URL:** `{{base_url}}/orders/non-existent-id`
- **Headers:**
    - Content-Type: application/json
    - Authorization: {{token}}
- **Body (raw, JSON):**
```json
{
  "status": "shipped"
}
```
- **Expected Response:**
    - **Status:** 404 Not Found
    - **Body:**
```json
{
  "error": "order not found: record not found"
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 404", function () {
    pm.response.to.have.status(404);
});
pm.test("Response contains error message", function () {
    pm.expect(pm.response.json().error).to.include("order not found");
});
```

### 3.9 GET /orders - List Orders (Success)

**Description:** List orders for the authenticated user.
- **Method:** GET
- **URL:** `{{base_url}}/orders?page=1&page_size=10`
- **Headers:**
    - Authorization: {{token}}
- **Body:** None
- **Expected Response:**
    - **Status:** 200 OK
    - **Body:**
```json
{
  "orders": [
    {
      "id": "{{order_id}}",
      "user_id": "{{token}}",
      "items": [
        {
          "product_id": "{{product_id}}",
          "quantity": 2
        }
      ],
      "status": "pending",
      "total": 999.98
    }
  ],
  "total": 1
}
```
- **Postman Tests:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
pm.test("Response contains orders array and total", function () {
    var jsonData = pm.response.