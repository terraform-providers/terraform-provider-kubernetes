package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	api "k8s.io/api/storage/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAccKubernetesCSIDriver_basic(t *testing.T) {
	var conf api.CSIDriver
	resourceName := "kubernetes_csi_driver.test"
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesCSIDriverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesCSIDriverBasicConfig(name, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesCSIDriverExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.attach_required", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKubernetesCSIDriverBasicConfig(name, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesCSIDriverExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.attach_required", "false"),
				),
			},
			{
				Config: testAccKubernetesCSIDriverBasicConfig(name, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesCSIDriverExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "spec.0.attach_required", "true"),
				),
			},
		},
	})
}

func testAccCheckKubernetesCSIDriverDestroy(s *terraform.State) error {
	conn, err := testAccProvider.Meta().(KubeClientsets).MainClientset()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_csi_driver" {
			continue
		}
		name := rs.Primary.ID
		resp, err := conn.StorageV1beta1().CSIDrivers().Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("CSIDriver still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesCSIDriverExists(n string, obj *api.CSIDriver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn, err := testAccProvider.Meta().(KubeClientsets).MainClientset()
		if err != nil {
			return err
		}

		name := rs.Primary.ID
		out, err := conn.StorageV1beta1().CSIDrivers().Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesCSIDriverBasicConfig(name string, attached bool) string {
	return fmt.Sprintf(`
resource "kubernetes_csi_driver" "test" {
  metadata {
    name = %[1]q
  }

  spec {
    attach_required   = %[2]t
  }
}
`, name, attached)
}
