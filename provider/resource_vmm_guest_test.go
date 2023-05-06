package provider

import (
	"fmt"
	"testing"

	"github.com/arnouthoebreckx/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSynologyVMMGuest(t *testing.T) {
	guestName := "terraform-guest"
	storageName := "synology - VM Storage 1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSynologyVMMGuestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSynologyVMMGuestConfig(guestName, storageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSynologyVMMGuestExists(guestName),
				),
			},
		},
	})
}

func testAccCheckSynologyVMMGuestDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(client.SynologyClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "synology_vmm_guest" {
			continue
		}

		guestName := rs.Primary.Attributes["guest_name"]

		err := client.DeleteGuest(guestName)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckSynologyVMMGuestConfig(guestName, storageName string) string {
	return fmt.Sprintf(`
		resource "synology_vmm_guest" "test_guest" {
			guest_name   = "%s"
			storage_name = "%s"
			description  = "Virtual Machine ACC test"
			autorun      = 2
			vram_size    = 1200

			vnics {
				network_name = "default"
			}

			vdisks {
				create_type = 0
				vdisk_size  = 20480
			}

			poweron = true
		}
	`, guestName, storageName)
}

func testAccCheckSynologyVMMGuestExists(guestName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["synology_vmm_guest.test_guest"]

		if !ok {
			return fmt.Errorf("Not found: synology_vmm_guest.test_guest")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set")
		}

		if rs.Primary.Attributes["guest_name"] != guestName {
			return fmt.Errorf("Guest name doesn't match")
		}

		return nil
	}
}
