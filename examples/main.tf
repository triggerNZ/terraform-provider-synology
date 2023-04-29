terraform {
  required_providers {
    synology = {
      version = "0.1"
      source  = "github.com/arnouthoebreckx/synology"
    }
  }
}

provider "synology" {
  url = "http://mylocation:5000"
  username = "myuser"
  password = "mypassword"
  # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}

# resource "synology_file" "hello-world" {
#   filename = "/home/downloaded/hello-world.txt"
#   content  = "Hello World"
# }

# resource "synology_file" "hello-world-from-file" {
#   filename = "/home/downloaded/hello-world-ref.txt"
#   content  = file("./hello-world.txt")
# }

# resource "synology_folder" "my-folder" {
#   path = "/home/downloaded/sample-folder"
# }

resource "synology_vmm_guest" "my-guest" {
  guest_name = "terraform-guest"
  storage_name = "storage1"
  vnics {
    network_name = "default"
  }
  vdisks {
    create_type = 0
    vdisk_size = 10240
  }
}
