package provider

import (
	"context"
    "strconv"
    "time"
    "log"

	"github.com/arnouthoebreckx/terraform-provider-synology/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func guestItem() *schema.Resource {
	return &schema.Resource{
		CreateContext: guestCreateItem,
		ReadContext:   guestReadItem,
        UpdateContext: guestUpdateItem,
        DeleteContext: guestDeleteItem,
		Schema: map[string]*schema.Schema{
			"guest_name": {
				Type:		schema.TypeString,
				Required:	true,
				Description: "The guest name",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The name of storage where the guest resides. Note: At least storage_id or storage_name should be given.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The id of storage where the guest resides. Note: At least storage_id or storage_name should be given.",
			},
            "vnics": {
                Type:       schema.TypeList,
                Required:   true,
                ForceNew:   true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "mac": {
                            Type: schema.TypeString,
                            Optional: true,
                            Description: "Optional. MAC address. If not specified, a MAC address of this vNIC will be randomly generated.",
                        },
                        "network_id": {
                            Type: schema.TypeString,
                            Optional: true,
                            Description: "Optional. Connected network group id. At least network_id or network_name should be given. Note: network_id can be an empty string to represent not being connected.",
                        },
                        "network_name": {
                            Type: schema.TypeString,
                            Optional: true,
                            Description: "Optional. Connected network group name. At least network_id or network_name should be given.",
                        },
                    },
                },
            },
            "vdisks": {
                Type:       schema.TypeList,
                Required:   true,
                ForceNew:   true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "create_type": {
                            Type: schema.TypeInt,
                            Required: true,
                            Description: "0: Create an empty vDisk, 1: Clone an existing image",
                        },
                        "vdisk_size": {
                            Type: schema.TypeInt,
                            Optional: true,
                            Description: "Optional. If create_type is 0, this field must be set. The created vDisk size in MB.",
                        },
                        "image_id": {
                            Type: schema.TypeString,
                            Optional: true,
                            Description: "Optional. If create_type is 1, at least image_id or image_name should be given. The id of the image that is to be cloned. Note: Image type should be disk.",
                        },
                        "image_name": {
                            Type: schema.TypeString,
                            Optional: true,
                            Description: "Optional. If create_type is 1, at least image_id or image_name should be given. The name of the image that is to be cloned. Note: Image type should be disk.",
                        },
                    },
                },
            },
		},
	}
}

func guestCreateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	name := d.Get("guest_name").(string)
    storage_id := d.Get("storage_id").(string)
    storage_name := d.Get("storage_name").(string)
    validateIdName(storage_id, storage_name)

    vnics := removeEmptyEntries(d.Get("vnics").([]interface{}))
	vdisks := removeEmptyEntries(d.Get("vdisks").([]interface{}))
    validateIdName(vnics, "network_id", "network_name")
    validateIdName(vdisks, "image_id", "image_name")

	service := GuestService{synologyClient: client}
	err := service.Create(name, storage_id, storage_name, vnics, vdisks)
	if err != nil {
		return diag.FromErr(err)
	}
	guestReadItem(ctx, d, m)
	return diags
}

func guestReadItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
    service := GuestService{synologyClient: client}
	name := d.Get("guest_name").(string)

    content, err := service.Read(name)
	if err != nil {
		return diag.FromErr(err)
	}

    log.Println("Read VMM Guest Content " + string(content))
	d.Set("guest_name", name)

    d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func guestUpdateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

	client := m.(client.SynologyClient)
    service := GuestService{synologyClient: client}

    if d.HasChange("guest_name") {
        name, new_name := d.GetChange("guest_name")
        err := service.Update(name.(string), new_name.(string))
        log.Println(err)
    }

    return diags
}

func guestDeleteItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

	client := m.(client.SynologyClient)
    service := GuestService{synologyClient: client}
    name := d.Get("guest_name").(string)
    
    err := service.Delete(name)

	if err != nil {
		return diag.FromErr(err)
	}

    return diags
}