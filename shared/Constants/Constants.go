package Constants

import "os"

const CORS_ALLOWED_ORIGINS = "*"
const DEFAULT_SERVER_PORT = "8090"
const SWAGGER_EDITOR_URL="https://editor.swagger.io/"

//env injected constants
var REDIS_CONNECTION_STRING = os.Getenv("REDIS_HOST")
var REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")