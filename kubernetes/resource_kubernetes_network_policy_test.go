package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	api "k8s.io/client-go/pkg/apis/networking/v1"
)

func TestAccKubernetesNetworkPolicy_basic(t *testing.T) {
	var conf api.NetworkPolicy
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_network_policy.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesNetworkPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesNetworkPolicyConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesNetworkPolicyExists("kubernetes_network_policy.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one"}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.TestLabelThree", "three"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.TestLabelFour", "four"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelThree": "three", "TestLabelFour": "four"}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.#", "0"),
				),
			},
			{
				Config: testAccKubernetesNetworkPolicyConfig_metaModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesNetworkPolicyExists("kubernetes_network_policy.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.TestAnnotationTwo", "two"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one", "TestAnnotationTwo": "two"}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.TestLabelThree", "three"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelTwo": "two", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.#", "0"),
				),
			},
			{
				Config: testAccKubernetesNetworkPolicyConfig_specModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesNetworkPolicyExists("kubernetes_network_policy.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.%", "0"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.%", "0"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.key", "name"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.operator", "In"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.values.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.values.1742479128", "webfront"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.values.2902841359", "api"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.0.port", "http"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.1.port", "8125"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.1.protocol", "UDP"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.0.namespace_selector.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.0.namespace_selector.0.match_labels.name", "default"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.0.pod_selector.#", "0"),
				),
			},
			{
				Config: testAccKubernetesNetworkPolicyConfig_specModified2(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesNetworkPolicyExists("kubernetes_network_policy.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.annotations.%", "0"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.labels.%", "0"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_network_policy.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.key", "name"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.operator", "In"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.values.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.values.1742479128", "webfront"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.pod_selector.0.match_expressions.0.values.2902841359", "api"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.0.port", "http"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.1.port", "statsd"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.ports.1.protocol", "UDP"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.0.namespace_selector.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.0.namespace_selector.0.match_labels.name", "default"),
					resource.TestCheckResourceAttr("kubernetes_network_policy.test", "spec.0.ingress.0.from.0.pod_selector.#", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesNetworkPolicy_importBasic(t *testing.T) {
	resourceName := "kubernetes_network_policy.test"
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesNetworkPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesNetworkPolicyConfig_basic(name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKubernetesNetworkPolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*kubernetes.Clientset)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_network_policy" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.NetworkingV1().NetworkPolicies(namespace).Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Namespace == namespace && resp.Name == name {
				return fmt.Errorf("Network Policy still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesNetworkPolicyExists(n string, obj *api.NetworkPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		out, err := conn.NetworkingV1().NetworkPolicies(namespace).Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesNetworkPolicyConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_network_policy" "test" {
  metadata {
    name      = "%s"
    namespace = "default"

    annotations {
      TestAnnotationOne = "one"
    }

    labels {
      TestLabelOne   = "one"
      TestLabelThree = "three"
      TestLabelFour  = "four"
    }
  }

  spec {
    pod_selector {}
  }
}
`, name)
}

func testAccKubernetesNetworkPolicyConfig_metaModified(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_network_policy" "test" {
  metadata {
    name      = "%s"
    namespace = "default"

    annotations {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }
  }

  spec {
    pod_selector = {}
    ingress      = []
  }
}
`, name)
}

func testAccKubernetesNetworkPolicyConfig_specModified(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_network_policy" "test" {
  metadata {
    name      = "%s"
    namespace = "default"
  }

  spec {
    pod_selector {
      match_expressions {
        key      = "name"
        operator = "In"
        values   = ["webfront", "api"]
      }
    }

    ingress = [
      {
        ports = [
          {
            port     = "http"
          },
          {
            port     = "8125"
            protocol = "UDP"
          },
        ]

        from = [
          {
            namespace_selector {
              match_labels = {
                name = "default"
              }
            }
          },
        ]
      },
    ]
  }
}
	`, name)
}

func testAccKubernetesNetworkPolicyConfig_specModified2(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_network_policy" "test" {
  metadata {
    name      = "%s"
    namespace = "default"
  }

  spec {
    pod_selector {
      match_expressions {
        key      = "name"
        operator = "In"
        values   = ["webfront", "api"]
      }
    }

    ingress = [
      {
        ports = [
          {
            port     = "http"
            protocol = "TCP"
          },
          {
            port     = "statsd"
            protocol = "UDP"
          },
        ]

        from = [
          {
            namespace_selector {
              match_labels = {
                name = "default"
              }
            }
          },
        ]
      },
    ]
  }
}
	`, name)
}
