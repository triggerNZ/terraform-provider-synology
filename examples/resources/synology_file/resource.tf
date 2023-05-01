terraform {
  required_providers {
    synology = {
      version = "0.1"
      source  = "github.com/arnouthoebreckx/synology"
    }
  }
}

provider "synology" {
  url      = "<SYNOLOGY_ADDRESS>"
  username = "<SYNOLOGY_USERNAME>"
  password = "<SYNOLOGY_PASSWORD>"
  # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}

resource "synology_file" "hello-world" {
  filename = "/home/downloaded/hello-world.txt"
  content  = "Hello World"
}

resource "synology_file" "hello-world-from-file" {
  filename = "/home/downloaded/hello-world-ref.txt"
  content  = file("./hello-world.txt")
}
