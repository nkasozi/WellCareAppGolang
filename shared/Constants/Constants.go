package Constants

import "os"

const CORS_ALLOWED_ORIGINS = "*"
const DEFAULT_SERVER_PORT = "8090"
const SWAGGER_EDITOR_URL="https://editor.swagger.io/"
const GOOGLE_CLOUD_PROJECT_ID = "reconcilercore"
const PUSH_NOTIFICATION_URL = "https://badadc09b6648dee9da7515d55e0ec68.m.pipedream.net"
const MESSAGE_ACK_DEALINE_IN_SECS = 60
const ENABLE_MESSAGE_ORDERING = true

//env injected constants
var REDIS_CONNECTION_STRING = os.Getenv("REDIS_HOST")
var REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")