# Go Clean Architecture

Proyek ini adalah implementasi dari arsitektur bersih (clean architecture) menggunakan bahasa pemrograman Go. Proyek ini menggunakan beberapa library seperti Fiber untuk web framework, GORM untuk ORM, dan Viper untuk manajemen konfigurasi.

## Arsitektur Proyek

Proyek ini diorganisir dengan beberapa lapisan utama:

1. **cmd**: Berisi entry point dari aplikasi.
2. **internal/config**: Berisi konfigurasi aplikasi seperti database, logger, validator, dan framework web.
3. **internal/delivery**: Berisi controller dan middleware untuk menangani HTTP request dan response.
4. **internal/entity**: Berisi definisi dari entitas-entitas yang digunakan dalam aplikasi.
5. **internal/helper**: Berisi helper functions yang digunakan di berbagai bagian aplikasi.
6. **internal/model**: Berisi definisi dari model-model yang digunakan untuk request dan response.
7. **internal/repository**: Berisi implementasi dari repository pattern untuk mengakses data.
8. **internal/usecase**: Berisi implementasi dari use case yang merupakan logika bisnis dari aplikasi.

## Cara Menjalankan Aplikasi

### Prasyarat

Pastikan Anda telah menginstal Go dan memiliki akses ke database MySQL.

### Langkah-langkah

1. Clone repository ini:
    ```sh
    git clone https://github.com/abdisetiakawan/go-clean-arch.git
    cd go-clean-arch
    ```

2. Buat file [config.json](http://_vscodecontentref_/0) di root directory dengan konfigurasi berikut:
    ```json
    {
        "app": {
            "name": "Go Clean Architecture"
        },
        "web": {
            "port": 8080,
            "prefork": false
        },
        "database": {
            "username": "your_db_username",
            "password": "your_db_password",
            "host": "localhost",
            "port": 3306,
            "name": "your_db_name",
            "pool": {
                "idle": 10,
                "max": 100,
                "lifetime": 300
            }
        },
        "log": {
            "level": 4
        },
        "credentials": {
            "accesssecret": "your_access_secret",
            "refreshsecret": "your_refresh_secret"
        }
    }
    ```

3. Jalankan migrasi database:
    ```sh
    go run cmd/migrate/main.go
    ```

4. Jalankan aplikasi:
    ```sh
    go run cmd/web/main.go
    ```

Aplikasi akan berjalan di [http://localhost:8080](http://_vscodecontentref_/1).

## Request dan Response

### User

#### Register User
- **Endpoint**: `POST /api/users`
- **Request Body**:
    ```json
    {
        "name": "John Doe",
        "email": "john.doe@example.com",
        "password": "password123"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "name": "John Doe",
            "email": "john.doe@example.com",
            "access_token": "access_token",
            "refresh_token": "refresh_token"
        }
    }
    ```

#### Login User
- **Endpoint**: `POST /api/users/_login`
- **Request Body**:
    ```json
    {
        "email": "john.doe@example.com",
        "password": "password123"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "name": "John Doe",
            "email": "john.doe@example.com",
            "access_token": "access_token",
            "refresh_token": "refresh_token"
        }
    }
    ```

#### Get Current User
- **Endpoint**: `GET /api/users/_current`
- **Response**:
    ```json
    {
        "data": {
            "name": "John Doe",
            "email": "john.doe@example.com"
        }
    }
    ```

#### Update User
- **Endpoint**: `PATCH /api/users/_current`
- **Request Body**:
    ```json
    {
        "name": "John Doe Updated",
        "password": "newpassword123"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "name": "John Doe Updated",
            "email": "john.doe@example.com"
        }
    }
    ```

### Task

#### Create Task
- **Endpoint**: `POST /api/tasks`
- **Request Body**:
    ```json
    {
        "title": "New Task",
        "description": "Task description",
        "status": "pending",
        "due_date": "2023-12-31"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "title": "New Task",
            "description": "Task description",
            "status": "pending",
            "due_date": "2023-12-31"
        }
    }
    ```

#### List Tasks
- **Endpoint**: `GET /api/tasks`
- **Response**:
    ```json
    {
        "data": [
            {
                "id": 1,
                "title": "New Task",
                "description": "Task description",
                "status": "pending",
                "due_date": "2023-12-31"
            }
        ],
        "paging": {
            "page": 1,
            "size": 10,
            "total_item": 1,
            "total_page": 1
        }
    }
    ```

#### Get Task
- **Endpoint**: `GET /api/tasks/:taskId`
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "title": "New Task",
            "description": "Task description",
            "status": "pending",
            "due_date": "2023-12-31"
        }
    }
    ```

#### Update Task
- **Endpoint**: `PUT /api/tasks/:taskId`
- **Request Body**:
    ```json
    {
        "title": "Updated Task",
        "description": "Updated description",
        "status": "in_progress",
        "due_date": "2023-12-31"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "title": "Updated Task",
            "description": "Updated description",
            "status": "in_progress",
            "due_date": "2023-12-31"
        }
    }
    ```

#### Delete Task
- **Endpoint**: `DELETE /api/tasks/:taskId`
- **Response**:
    ```json
    {
        "data": true
    }
    ```

### Tag

#### Create Tag
- **Endpoint**: `POST /api/tags`
- **Request Body**:
    ```json
    {
        "name": "New Tag"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "name": "New Tag"
        }
    }
    ```

#### List Tags
- **Endpoint**: `GET /api/tags`
- **Response**:
    ```json
    {
        "data": [
            {
                "id": 1,
                "name": "New Tag"
            }
        ],
        "paging": {
            "page": 1,
            "size": 10,
            "total_item": 1,
            "total_page": 1
        }
    }
    ```

#### Get Tag
- **Endpoint**: `GET /api/tags/:tagId`
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "name": "New Tag"
        }
    }
    ```

#### Update Tag
- **Endpoint**: `PUT /api/tags/:tagId`
- **Request Body**:
    ```json
    {
        "name": "Updated Tag"
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "name": "Updated Tag"
        }
    }
    ```

#### Delete Tag
- **Endpoint**: `DELETE /api/tags/:tagId`
- **Response**:
    ```json
    {
        "data": true
    }
    ```

### Task Tag

#### Create Task Tag
- **Endpoint**: `POST /api/tasks/:taskId/tags`
- **Request Body**:
    ```json
    {
        "tag_id": 1
    }
    ```
- **Response**:
    ```json
    {
        "data": {
            "id": 1,
            "taskId": 1,
            "tag_id": 1
        }
    }
    ```

#### List Task Tags
- **Endpoint**: `GET /api/taskswithtags`
- **Response**:
    ```json
    {
        "data": [
            {
                "id": 1,
                "title": "New Task",
                "description": "Task description",
                "status": "pending",
                "due_date": "2023-12-31",
                "tag_id": 1
            }
        ],
        "paging": {
            "page": 1,
            "size": 10,
            "total_item": 1,
            "total_page": 1
        }
    }
    ```

#### List Tasks by Tag ID
- **Endpoint**: `GET /api/tags/:tagId/tasks`
- **Response**:
    ```json
    {
        "data": [
            {
                "id": 1,
                "title": "New Task",
                "description": "Task description",
                "status": "pending",
                "due_date": "2023-12-31",
                "tag_id": 1
            }
        ],
        "paging": {
            "page": 1,
            "size": 10,
            "total_item": 1,
            "total_page": 1
        }
    }
    ```

#### Delete Task Tag
- **Endpoint**: `DELETE /api/tasks/:taskId/tags/:tagId`
- **Response**:
    ```json
    {
        "data": true
    }
    ```
