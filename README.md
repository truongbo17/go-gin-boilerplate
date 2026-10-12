# Go Gin Boilerplate

Boilerplate Go dùng Gin, có sẵn REST API, GraphQL, JWT, phân quyền RBAC, MySQL migration, rate limit và worker xử lý tác vụ nền. Mã được chia theo feature để bạn có thể thay auth mẫu hoặc thêm nghiệp vụ riêng.

## Bắt đầu nhanh

### Cần chuẩn bị

- Go theo phiên bản khai báo trong `go.mod`
- Docker Compose để chạy MySQL; Redis chỉ cần khi dùng worker hoặc cấu hình cache Redis
- Python 3 để nhập mật khẩu quản trị mà không hiện trên màn hình ở bước dưới

### 1. Tạo cấu hình local

```sh
cp .env.example .env
openssl rand -hex 32
```

Điền chuỗi vừa tạo vào `JWT_SECRET` trong `.env`. Đặt mật khẩu cho `DB_PASS` và `MYSQL_ROOT_PASSWORD`.

### 2. Chạy MySQL

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mysql
```

Ứng dụng chạy trên máy của bạn, vì vậy giữ `DB_HOST=127.0.0.1` và `DB_PORT=3306` trong `.env`.

### 3. Tạo bảng và tài khoản quản trị

```sh
go run . migrate
python3 -c 'import getpass; print(getpass.getpass("Mật khẩu admin: "))' | go run . create_user \
  --username admin \
  --email admin@example.com \
  --password-stdin \
  --admin
```

`--admin` chỉ dành cho người vận hành đáng tin cậy. API không cho phép tạo, đổi tên hoặc xóa role `admin`.

### 4. Chạy API

```sh
go run . server
```

Mặc định API nghe tại `http://localhost:8000`. Kiểm tra nhanh:

```sh
curl http://localhost:8000/ping
```

Khi chạy `APP_ENV=local` hoặc `debug`, mở `http://localhost:8000/swagger/index.html` để xem REST API. Swagger và OpenAPI bị tắt trong môi trường `release`.

## Tính năng có sẵn

| Đường dẫn | Mô tả |
| --- | --- |
| `POST /api/v1/auth/register` | Tạo tài khoản |
| `POST /api/v1/auth/login` | Đăng nhập, nhận JWT |
| `POST /api/v1/auth/forgot-password` | Gửi yêu cầu đặt lại mật khẩu |
| `POST /api/v1/auth/reset-password` | Đặt mật khẩu mới bằng token |
| `/api/v1/rbac/*` | Quản lý user, role và permission |
| `POST /api/v1/graphql` | GraphQL có xác thực |
| `GET /api/v1/public/enum-options/:key` | Đọc danh mục enum công khai |
| `GET /api/v1/graphql/entity-options` | Danh sách entity có kiểm tra quyền |

Các route cần đăng nhập nhận header `Authorization: Bearer <token>`. Đổi mật khẩu sẽ làm mất hiệu lực các access token cũ của tài khoản đó.

## Worker và email mẫu

Worker cần Redis. Trong `.env`, đặt `CACHE_STORE=redis`, `REDIS_HOST=127.0.0.1`, `REDIS_PORT=6379`, sau đó chạy:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d redis
go run . worker
```

Worker hiện có ví dụ schedule in `schedule` mỗi 15 giây và các job email chào mừng/đặt lại mật khẩu.

Để thử email local, chạy thêm Mailpit:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mailpit
```

Đặt `MAIL_ENABLED=true`, `MAIL_HOST=127.0.0.1`, `MAIL_PORT=1025`, `MAIL_FROM=noreply@example.test`, `MAIL_TLS_MODE=none` và `PASSWORD_RESET_URL=http://localhost:3000/reset-password`. Chạy `server` và `worker` ở hai terminal. Xem thư tại `http://127.0.0.1:8025`; email không được gửi ra ngoài máy.

Ở production, dùng `MAIL_TLS_MODE=starttls` hoặc `implicit`, URL HTTPS và lưu thông tin SMTP trong secret manager. Queue có retry nhưng không bảo đảm gửi đúng một lần; email có thể được gửi lặp khi kết nối SMTP lỗi ở thời điểm xác nhận.

## Cấu hình thường dùng

Các cấu hình được đọc từ `.env` hoặc biến môi trường. Giá trị quan trọng:

- `APP_ENV`: `local`, `debug` hoặc `release`
- `APP_PORT`: cổng HTTP, mặc định `8000`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_DATABASE`: kết nối MySQL
- `CACHE_STORE`: `local` cho một instance, `redis` để chia sẻ cache/rate limit giữa nhiều instance
- `JWT_SECRET`: khóa riêng, tối thiểu 32 ký tự
- `MAIL_ENABLED`: bật hoặc tắt job email
- `TRACER_ENABLE`: bật xuất OpenTelemetry

Danh sách biến mẫu nằm trong [.env.example](.env.example). Không commit file `.env` hoặc secret thật.

## Cấu trúc source

- `cmd/`: lệnh `server`, `worker`, `migrate`, `create_user`, `version` và vòng đời tài nguyên
- `config/`: đọc và kiểm tra cấu hình
- `internal/app/core/<feature>/`: nghiệp vụ và kiểu dữ liệu độc lập HTTP
- `internal/app/v1/<feature>/`: route, controller, request và response của REST
- `internal/app/worker/<feature>/`: handler job; handler gọi về service nghiệp vụ
- `internal/repository/<feature>/`: truy vấn và giao dịch database
- `internal/infra/`: client cho database, Redis, Kafka, MongoDB, ClickHouse, S3, SMTP, tracing và worker
- `internal/routes/`: cấu hình Gin, middleware chung và nơi ghép các feature HTTP

Khi thêm một feature HTTP, tạo module và `RegisterRoutes` trong `internal/app/v1/<feature>`, sau đó khởi tạo và gắn route tại `internal/routes/main.go`. Đặt rate limit, xác thực, quyền, validate và controller ngay cạnh route để dễ đọc thứ tự xử lý.

Xem [docs/architecture.md](docs/architecture.md) về luồng request, worker và hạ tầng; [docs/go-style.md](docs/go-style.md) về quy ước code.

## Lệnh phát triển

```sh
make fmt
make vet
make build
go run . version
```

Air có thể tự build và khởi động lại server khi sửa code:

```sh
go install github.com/air-verse/air@v1.67.4
air -c .air.toml
```

## Mặc định bảo mật

- Đặt `JWT_SECRET` riêng cho từng môi trường; không dùng khóa ví dụ.
- Chỉ khai báo CORS origin của frontend cần truy cập.
- Khi đặt sau reverse proxy, chỉ tin cậy IP/CIDR proxy cụ thể trong `APP_TRUSTED_PROXIES`.
- Bật Redis cache khi chạy nhiều instance để chia sẻ token revoke và rate limit.
- Đặt TLS ở reverse proxy và không mở trực tiếp cổng ứng dụng ra Internet nếu proxy quản lý truy cập.
- Tạo tài khoản qua CLI; dùng `--password-stdin` để mật khẩu không hiện trong lịch sử shell.
