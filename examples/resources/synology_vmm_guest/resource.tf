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

resource "synology_vmm_guest" "my-guest" {
  auto_run     = 2
  poweron      = true
  guest_name   = "terraform-guest"
  description  = "Virtual machine setup with terraform"
  storage_name = "synology - VM Storage 1"
  vram_size    = 1024
  vnics {
    network_name = "default"
  }
  vdisks {
    create_type = 0
    vdisk_size  = 10240
  }
}
