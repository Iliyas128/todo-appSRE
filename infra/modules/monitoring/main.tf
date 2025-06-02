terraform {
  required_providers {
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

resource "docker_image" "prometheus" {
  name = "prom/prometheus"
}

resource "docker_container" "prometheus" {
  name  = "prometheus"
  image = docker_image.prometheus.name
  networks_advanced {
    name = "todo_net"
  }
  ports {
    internal = 9090
    external = 9090
  }
}

resource "docker_image" "grafana" {
  name = "grafana/grafana"
}

resource "docker_container" "grafana" {
  name  = "grafana"
  image = docker_image.grafana.name
  networks_advanced {
    name = "todo_net"
  }
  ports {
    internal = 3000
    external = 3000
  }
}
