resource "scaleway_object_bucket" "opdb-front" {
  name = "opdb-front"
}

resource "scaleway_object_bucket" "opdb-data" {
  region = "fr-par"
  name = "opdb-data"
}

resource "scaleway_object_bucket_acl" "opdb-data" {
  bucket = "opdb-data"
  acl = "private"
}