package crawler

// Regex
const (
    URL_REGEX = `(http|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`
)

// Errors
const (
    PARSE_ERROR = "Failed to parse"
)

// User agents
const (
    USER_AGENT = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36`
)


// Store keys
type StoreKey string

const (
    StoreKeyUncrawled StoreKey = "uncrawled"
    StoreKeyAll = "all"
)