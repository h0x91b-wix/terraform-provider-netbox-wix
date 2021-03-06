package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/h0x91b-wix/go-netbox/netbox/client/ipam"
	"github.com/h0x91b-wix/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceNetboxRegionalInternetRegistry is the core Terraform resource structure for the netbox_regional_internet_registry resource.
func resourceNetboxRegionalInternetRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxRegionalInternetRegistryCreate,
		Read:   resourceNetboxRegionalInternetRegistryRead,
		Update: resourceNetboxRegionalInternetRegistryUpdate,
		Delete: resourceNetboxRegionalInternetRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"rir_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

// resourceNetboxRegionalInternetRegistryCreate creates a new RIR in Netbox.
func resourceNetboxRegionalInternetRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	isPrivate := d.Get("is_private").(bool)

	var parm = ipam.NewIpamRirsCreateParams().WithData(
		&models.RIR{
			Slug:      &slug,
			Name:      &name,
			IsPrivate: isPrivate,
		},
	)

	log.Debugf("Executing IpamRirsCreate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamRirsCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamRirsCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(fmt.Sprintf("ipam/rir/%v", out.Payload.ID))
	d.Set("rir_id", out.Payload.ID)

	log.Debugf("Done Executing IpamRirsCreate: %v", out)

	return nil
}

// resourceNetboxRegionalInternetRegistryUpdate applies updates to a RIR by ID when deltas are detected by Terraform.
func resourceNetboxRegionalInternetRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	//terraformID := d.Id()
	netboxID := int64(d.Get("rir_id").(int))
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	isPrivate := d.Get("is_private").(bool)

	var parm = ipam.NewIpamRirsUpdateParams().
		WithID(netboxID).
		WithData(
			&models.RIR{
				Slug:      &slug,
				Name:      &name,
				IsPrivate: isPrivate,
			},
		)

	log.Debugf("Executing IpamRirsUpdate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamRirsUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamRirsUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IpamRirsUpdate: %v", out)

	return nil
}

// resourceNetboxRegionalInternetRegistryRead reads an existing RIR by ID.
func resourceNetboxRegionalInternetRegistryRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	//terraformID := d.Id()
	netboxID := int64(d.Get("rir_id").(int))

	var readParams = ipam.NewIpamRirsReadParams().WithID(netboxID)

	readRirResult, err := netboxClient.Ipam.IpamRirsRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching RIR ID # %d from Netbox = %v", netboxID, err)
		return err
	}

	log.Debugf("Read RIR %d = %v", netboxID, readRirResult.Payload)

	d.Set("name", readRirResult.Payload.Name)
	d.Set("slug", readRirResult.Payload.Slug)
	d.Set("is_private", readRirResult.Payload.IsPrivate)
	d.Set("rir_id", readRirResult.Payload.ID)

	return nil
}

// resourceNetboxRegionalInternetRegistryDelete deletes an existing RIR by ID.
func resourceNetboxRegionalInternetRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting RIR: %v\n", d)

	netboxID := int64(d.Get("rir_id").(int))

	var deleteParameters = ipam.NewIpamRirsDeleteParams().WithID(netboxID)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.Ipam.IpamRirsDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamRirsDelete: %v", err)
	}

	log.Debugf("Done Executing IpamRirsDelete: %v", out)

	return nil
}
