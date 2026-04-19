# 01 — Clean Code: Production Go Engineering Standards

---

## Function Design Rules

### Single Responsibility, Hard Limit

A function does one thing. If you need "and" in the description, split it.

```go
// BAD: fetch, parse, and persist in one shot
func ProcessUser(id string) error {
    resp, err := http.Get("/users/" + id)
    // ... parse, validate, write to DB
}

// GOOD: separated concerns, each independently testable
func FetchUser(ctx context.Context, id string) (*UserDTO, error) { ... }
func ValidateUser(u *UserDTO) error { ... }
func PersistUser(ctx context.Context, db *sql.DB, u *User) error { ... }
```

### Return Early, Never Nest

Nesting is a readability tax. Pay it up front with guard clauses.

```go
// BAD
func Handle(r *Request) (*Response, error) {
    if r != nil {
        if r.Body != nil {
            if r.Method == "POST" {
                // actual logic buried 3 levels deep
            }
        }
    }
    return nil, nil
}

// GOOD
func Handle(r *Request) (*Response, error) {
    if r == nil {
        return nil, ErrNilRequest
    }
    if r.Body == nil {
        return nil, ErrEmptyBody
    }
    if r.Method != "POST" {
        return nil, ErrMethodNotAllowed
    }
    // actual logic at top level
}
```

### Function Signatures: Context First, Error Last

```go
// Canonical Go function shape for I/O-bound work
func (s *Service) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error)

// Options pattern for complex construction — never > 5 params
type ServerOption func(*serverConfig)

func WithTimeout(d time.Duration) ServerOption {
    return func(c *serverConfig) { c.timeout = d }
}

func NewServer(addr string, opts ...ServerOption) *Server {
    cfg := defaultConfig()
    for _, o := range opts {
        o(cfg)
    }
    return &Server{addr: addr, cfg: cfg}
}
```

### Zero-Value Usability

Design structs so zero value is valid and useful. Avoid required constructor footguns.

```go
// BAD: zero value is broken
type Cache struct {
    mu   sync.Mutex
    data map[string]any // nil map panics on write
}

// GOOD: lazy init, zero value is safe
type Cache struct {
    mu   sync.Mutex
    data map[string]any
}

func (c *Cache) Set(k string, v any) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.data == nil {
        c.data = make(map[string]any)
    }
    c.data[k] = v
}
```

---

## Naming Conventions

### Rules That Aren't Negotiable

| Context | Rule | Example |
|---|---|---|
| Exported types | PascalCase, noun | `OrderRepository` |
| Unexported vars | camelCase, short in scope | `cfg`, `req`, `mu` |
| Interfaces | `-er` suffix where applicable | `Reader`, `Closer`, `EventEmitter` |
| Error vars | `Err` prefix | `ErrNotFound`, `ErrTimeout` |
| Error types | `-Error` suffix | `ValidationError`, `NetworkError` |
| Booleans | `is/has/can` prefix | `isReady`, `hasExpired` |
| Receivers | 1–2 letter abbreviation of type | `(s *Server)`, `(r *Repository)` |

### Acronyms: All Caps or All Lower

```go
// CORRECT
type HTTPClient struct{}
type userID string
func parseURL(s string) (*url.URL, error)
var sqlDB *sql.DB

// WRONG
type HttpClient struct{}
type userId string
func parseUrl(s string) (*url.URL, error)
```

### Don't Stutter the Package Name

```go
// Package: auth
// BAD
auth.AuthService
auth.AuthMiddleware

// GOOD
auth.Service
auth.Middleware
```

### Loop Variables: Short is Correct

```go
for i, v := range items { ... }       // i, v are fine
for _, order := range orders { ... }  // order is fine
for k, v := range m { ... }           // k, v are fine
```

---

## Package Architecture

### Flat Is Better Than Nested

The standard Go project layout for a service:

```
myservice/
├── cmd/
│   └── server/
│       └── main.go          // wire everything together here
├── internal/
│   ├── handler/             // HTTP/gRPC handlers
│   ├── service/             // business logic
│   ├── repository/          // data access
│   └── domain/              // core types, no imports from above layers
├── pkg/                     // exported, reusable across services
│   ├── middleware/
│   └── telemetry/
└── config/
    └── config.go
```

### The Dependency Rule

```
handler → service → repository → domain
```

`domain` imports nothing internal. `repository` imports `domain`. `service` imports `repository` + `domain`. `handler` imports `service`. Reverse imports are a design error — use interfaces to invert.

### Internal Packages Are Not Optional

Use `internal/` to enforce boundaries. The Go toolchain enforces it at compile time. Any code outside the module cannot import `internal/`.

```go
// internal/repository/order.go
package repository

type OrderStore interface {
    FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
    Save(ctx context.Context, o *domain.Order) error
}

type postgresOrderStore struct {
    db *pgxpool.Pool
}

// Only exported via the interface. Concrete type stays internal.
func NewOrderStore(db *pgxpool.Pool) OrderStore {
    return &postgresOrderStore{db: db}
}
```

### Interface Segregation: Small Interfaces Win

```go
// BAD: fat interface forces fake implementations in tests
type UserService interface {
    Create(ctx context.Context, req CreateUserReq) (*User, error)
    Update(ctx context.Context, id string, req UpdateUserReq) (*User, error)
    Delete(ctx context.Context, id string) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, id string) (*User, error)
    ListActive(ctx context.Context) ([]*User, error)
    // ... 10 more methods
}

// GOOD: consumers define what they need
type UserFinder interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserCreator interface {
    Create(ctx context.Context, req CreateUserReq) (*User, error)
}
```

---

## Composition Over Inheritance

Go has no inheritance. Embed for behavior composition, not type hierarchy.

```go
// Base logger behavior
type BaseLogger struct {
    level  slog.Level
    logger *slog.Logger
}

func (b *BaseLogger) Info(msg string, args ...any) {
    b.logger.Info(msg, args...)
}

// Compose: OrderService gets logging for free, no interface pollution
type OrderService struct {
    BaseLogger
    repo  repository.OrderStore
    cache *redis.Client
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderReq) (*Order, error) {
    s.Info("creating order", "user_id", req.UserID) // from BaseLogger
    // ...
}
```

### Embedding Interfaces for Partial Implementations

```go
// Test doubles without implementing the full interface
type mockStore struct {
    repository.OrderStore // embed to satisfy interface
    findByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Order, error)
}

func (m *mockStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
    return m.findByIDFn(ctx, id)
}
// Only override what the test needs. Other methods panic at runtime if called.
```

### Avoid Deep Embedding Chains

```go
// This becomes unreadable fast
type A struct{}
type B struct{ A }
type C struct{ B }
type D struct{ C }

// Prefer explicit delegation when the chain is > 2 deep
type D struct {
    c *C
}
func (d *D) DoThing() { d.c.DoThing() }
```

---

## Error Handling

### Sentinel Errors vs Error Types

```go
// Sentinel: for expected domain conditions, compared with ==  or errors.Is
var (
    ErrNotFound   = errors.New("not found")
    ErrConflict   = errors.New("conflict")
    ErrForbidden  = errors.New("forbidden")
)

// Error type: when you need to carry context
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// Caller can type-assert for rich info
var ve *ValidationError
if errors.As(err, &ve) {
    log.Printf("bad field: %s", ve.Field)
}
```

### Wrapping: fmt.Errorf + %w

```go
func (r *postgresOrderStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
    row := r.db.QueryRow(ctx, queryFindOrder, id)
    o := &domain.Order{}
    if err := row.Scan(&o.ID, &o.Status, &o.Total); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("order %s: %w", id, ErrNotFound)
        }
        return nil, fmt.Errorf("FindByID scan: %w", err)
    }
    return o, nil
}
```

### Never Discard Errors in Production Code

```go
// BAD — silently swallows failures
_ = cache.Set(ctx, key, val)

// GOOD — log non-critical, propagate critical
if err := cache.Set(ctx, key, val); err != nil {
    s.logger.Warn("cache write failed, continuing", "err", err, "key", key)
}
```

### Error Propagation Chain Rule

Wrap once with context at each layer boundary. Don't re-wrap the same message three times.

```go
// repository layer
return nil, fmt.Errorf("orderRepo.FindByID: %w", err)

// service layer
return nil, fmt.Errorf("OrderService.GetOrder: %w", err)

// handler layer — log the full chain, respond with sanitized message
if errors.Is(err, ErrNotFound) {
    http.Error(w, "order not found", http.StatusNotFound)
    return
}
s.logger.Error("get order failed", "err", err, "order_id", id)
http.Error(w, "internal error", http.StatusInternalServerError)
```

### panic Is Not Error Handling

Reserve `panic` for unrecoverable programmer errors (nil dereference in init, broken invariants). Never panic for runtime conditions like network failures or missing config values.

```go
// Acceptable panic: programming error, should never happen in correct code
func mustCompileRegex(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil {
        panic(fmt.Sprintf("invalid regex %q: %v", pattern, err))
    }
    return re
}

var emailRegex = mustCompileRegex(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
```
