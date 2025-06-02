terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

provider "docker" {}

# Network configuration
resource "docker_network" "todo_network" {
  name = "todo-network"
  driver = "bridge"
}

# Volume for persistent data
resource "docker_volume" "todo_data" {
  name = "todo-data"
}

# Todo App Service
resource "docker_container" "todo_app" {
  name  = "todo-app"
  image = "todo-app:latest"
  ports {
    internal = 8080
    external = 8080
  }
  networks_advanced {
    name = docker_network.todo_network.name
  }
  depends_on = [docker_network.todo_network]
}

# Prometheus Service
resource "docker_container" "prometheus" {
  name  = "prometheus"
  image = "prom/prometheus:latest"
  ports {
    internal = 9090
    external = 9090
  }
  volumes {
    host_path      = abspath("${path.module}/prometheus.yml")
    container_path = "/etc/prometheus/prometheus.yml"
  }
  networks_advanced {
    name = docker_network.todo_network.name
  }
  depends_on = [docker_network.todo_network]
}

# Grafana Service
resource "docker_container" "grafana" {
  name  = "grafana"
  image = "grafana/grafana:latest"
  ports {
    internal = 3000
    external = 3000
  }
  env = [
    "GF_SECURITY_ADMIN_USER=admin",
    "GF_SECURITY_ADMIN_PASSWORD=admin"
  ]
  networks_advanced {
    name = docker_network.todo_network.name
  }
  depends_on = [docker_network.todo_network]
}
