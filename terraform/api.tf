resource scaleway_container_namespace opdb-ns {
  name = "opdb"
  description = "Namespace managed by terraform"
}

resource scaleway_container opdb-api {
  name = "opdb-api"
  description = "API container deployed with terraform"
  namespace_id = scaleway_container_namespace.opdb-ns.id
  registry_image = "rg.fr-par.scw.cloud/opdb/api:v0.0.0-1"
  port = 8080
  cpu_limit = 70
  memory_limit = 128
  min_scale = 0
  max_scale = 1
  privacy = "public"
  protocol = "http1"
  deploy = true
  max_concurrency = 80
}

output "endpoint" {
  value = scaleway_container.opdb-api.domain_name
}