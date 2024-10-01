package proxy

const (
	// Version
	SOCKS_VERSION byte = 0x05

	// Supported authentication methods
	SOCKS_AUTH_METHOD_NO_REQUIRED           byte = 0x00
	SOCKS_AUTH_METHOD_NO_ACCEPTABLE_METHODS byte = 0xFF

	// Supported command
	SOCKS_CMD_CONNECT byte = 0x01

	// Reserved byte
	SOCKS_RESERVED_BYTE byte = 0x00

	// Supported address type
	SOCKS_ADDR_TYPE_IPV4 byte = 0x01
	SOCKS_ADDR_TYPE_FQDN byte = 0x03

	// Reply
	SOCKS_REPLY_SUCCEEDED                         byte = 0x00
	SOCKS_REPLY_GENERAL_SOCKS_SERVER_FAILURE      byte = 0x01
	SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET byte = 0x02
	SOCKS_REPLY_NETWORK_UNREACHABLE               byte = 0x03
	SOCKS_REPLY_HOST_UNREACHABLE                  byte = 0x04
	SOCKS_REPLY_CONNECTION_REFUSED                byte = 0x05
	SOCKS_REPLY_TTL_EXPIRED                       byte = 0x06
	SOCKS_REPLY_COMMAND_NOT_SUPPORTED             byte = 0x07
	SOCKS_REPLY_ADDRESS_TYPE_NOT_SUPPORTED        byte = 0x08
)
