terraform {
  required_providers {
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

variable "app_image" {
  type = string
}

resource "docker_image" "todo_app" {
  name         = var.app_image
  keep_locally = true
}

resource "docker_container" "todo_app" {
  name  = "todo_app"
  image = docker_image.todo_app.name
  networks_advanced {
    name = "todo_net"
  }
  ports {
    internal = 8080
    external = 8080
  }
}
