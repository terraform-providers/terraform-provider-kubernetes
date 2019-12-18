package kubernetes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceKubernetesPersistentVolumeClaim() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubernetesPersistentVolumeClaimRead,

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("persistent volume claim", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the desired characteristics of a volume requested by a pod author. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#persistentvolumeclaims",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_modes": {
							Type:        schema.TypeSet,
							Description: "A set of the desired access modes the volume should have. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#access-modes-1",
							Computed:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						"resources": {
							Type:        schema.TypeList,
							Description: "A list of the minimum resources the volume should have. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources",
							Computed:    true,
							ForceNew:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limits": {
										Type:        schema.TypeMap,
										Description: "Map describing the maximum amount of compute resources allowed. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
										Optional:    true,
										ForceNew:    true,
									},
									"requests": {
										Type:        schema.TypeMap,
										Description: "Map describing the minimum amount of compute resources required. If this is omitted for a container, it defaults to `limits` if that is explicitly specified, otherwise to an implementation-defined value. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
						"selector": {
							Type:        schema.TypeList,
							Description: "A label query over volumes to consider for binding.",
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(false),
							},
						},
						"volume_name": {
							Type:        schema.TypeString,
							Description: "The binding reference to the PersistentVolume backing this claim.",
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
						},
						"storage_class_name": {
							Type:        schema.TypeString,
							Description: "Name of the storage class requested by the claim",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKubernetesPersistentVolumeClaimRead(d *schema.ResourceData, meta interface{}) error {
	metadata := expandMetadata(d.Get("metadata").([]interface{}))

	om := meta_v1.ObjectMeta{
		Namespace: metadata.Namespace,
		Name:      metadata.Name,
	}
	d.SetId(buildId(om))

	return resourceKubernetesPersistentVolumeClaimRead(d, meta)
}
