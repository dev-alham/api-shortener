package config

const (
	/* config for database */
	PREFIX_DB_NAME = "shortener"
	PRINT_QUERY    = true

	/* jwt */
	EX_TIME_JWT = 1

	/* redis */
	DIR_CACHE_AUTH      = "AUTH"
	CACHE_DIR_SHORT_URL = "SHORT URL"
	CACHE_DIR_LONG_URL  = "LONG URL"
	CACHE_DIR_LIMIT     = "LIMIT"

	/* request */
	LIMIT_REQUEST_GET_DAY  = 10
	LIMIT_REQUEST_POST_DAY = 3
)
