ui = true
listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}
storage "file" {
  path = "/vault/data"
}
plugin_directory = "/lib/vault/plugin"
api_addr = "http://0.0.0.0:8200"