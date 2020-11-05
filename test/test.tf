terraform {
  required_providers {
    elasticsearch = {
      versions = ["0.1"]
      source   = "registry.terraform.io/estaldo/elasticsearch"
    }
  }
}

provider "elasticsearch" {
  url      = "http://localhost:9200"
  username = "elastic"
  password = "nope"
}

resource "elasticsearch_user" "test" {
  username ="username"
  password = "password"
  enabled = false
  email = "email@email.nope"
  full_name = "Full Name"
  roles = [
    "test"
  ]
  metadata = {
    "meta1": "value1",
    "meta2": "value2",
  }
}

resource "elasticsearch_role" "test" {
  name = "test"
  cluster = ["all"]
  indices {
    names = ["*"]
    privileges = ["all"]
    field_security {
      grant = ["*"]
    }
  }
  metadata = {
    "meta1": "value1",
    "meta2": "value2",
  }
  
  depends_on = [
    elasticsearch_user.test
  ]
}

resource "elasticsearch_api_key" "test" {
  name = "test2"
  expiration = "7d"
}
