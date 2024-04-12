variable "api_container" {
  type = string
}

resource scaleway_container_namespace opdb-ns {
  name = "opdb"
  description = "Namespace managed by terraform"
}

resource scaleway_container opdb-api {
  name = "opdb-api"
  description = "API container deployed with terraform"
  namespace_id = scaleway_container_namespace.opdb-ns.id
  registry_image = var.api_container
  port = 8080
  cpu_limit = 70
  memory_limit = 128
  min_scale = 0
  max_scale = 1
  privacy = "public"
  protocol = "h2c"
  http_option = "redirected"
  deploy = true
  max_concurrency = 80
}

output "endpoint" {
  value = scaleway_container.opdb-api.domain_name
}