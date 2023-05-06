package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSynologyVMMGuestDataSource(t *testing.T) {
	guestName := "terraform-guest"
	storageName := "synology - VM Storage 1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSynologyVMMGuestDataSourceConfig(guestName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "guest_name", guestName),
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "storage_name", storageName),
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "autorun", fmt.Sprint(2)),
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "description", "Virtual Machine ACC test"),
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "vram_size", fmt.Sprint(1200)),
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "vnics.0.network_name", "default"),
					resource.TestCheckResourceAttr("synology_vmm_guest.test_guest", "vdisks.0.vdisk_size", fmt.Sprint(20480)),
				),
			},
		},
	})
}

func testAccSynologyVMMGuestDataSourceConfig(guestName string) string {
	return fmt.Sprintf(`
		data "synology_vmm_guest" "my-guest" {
			guest_name = "%s"
		}
	`, guestName)
}
