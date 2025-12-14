# simple-router

Provides a simple router interface with middleware support and not much else.

## Exports

- Init
- SetNotFound
- SetRecovery
- AddMiddleware
- AppendMiddleware
- Add
- Get
- HeadOnly

<hr>

### `Init() *temp`
**Description:** Initializes a handle for the builder.

**Returns:**
- `*temp`

<hr>

### `(*temp) SetNotFound(handler http.HandlerFunc)`
**Description:** Sets which handler will be called when the method or URL path don't match.

**Args:**
- `handler`

<hr>

### `(*temp) SetRecovery(h func(any, http.ResponseWriter, *http.Request))`
**Description:** Sets which function will be called when recovering a panic.

**Args:**
- `h` - handler with an extra arg `err := recover()`

<hr>

### `(*temp) AddMiddleware(middleware func(http.HandlerFunc) http.HandlerFunc)`
**Description:** Creates and subsequently appends middleware that will be used at the very beginning of every request.

**Args:**
- `middleware`

<hr>

### `(*temp) AppendMiddleware(class int, middleware func(http.HandlerFunc) http.HandlerFunc)`
**Description:** Creates and subsequently appends middleware that will be used for every request added with the given class.

**Args:**
- `class` - nonnegative int
- `middleware`

<hr>

### `(*temp) Add(method string, class int, path string, handler http.HandlerFunc)`
**Description:** Adds a route.

**Args:**
- `method` - http standard string
- `class` - nonnegative int
- `path`
- `handler`

<hr>

### `(*temp) Get() http.HandlerFunc`
**Description:** Builds and returns the router.

**Returns:**
- `http.HandlerFunc`

<hr>

### `HeadOnly(handler http.HandlerFunc) http.HandlerFunc`
**Description:** Returns the input without the capability of writing the body.

**Args:**
- `handler`

**Returns:**
- `http.HandlerFunc`
