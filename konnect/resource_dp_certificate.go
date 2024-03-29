package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/csechrist/terraform-provider-konnect/konnect/client"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDPCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDPCertificateCreate,
		ReadContext:   resourceDPCertificateRead,
		DeleteContext: resourceDPCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, math.MaxInt),
				ForceNew:     true,
			},
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func fillDPCertificate(c *client.DataPlaneCertificate, d *schema.ResourceData) {
	c.Cert = d.Get("cert").(string)
	c.ControlPlaneId = d.Get("control_plane_id").(string)
}

func fillResourceDataFromDPCertificate(c *client.DataPlaneCertificate, d *schema.ResourceData) {
	d.Set("id", c.Id)
	d.Set("cert", c.Cert)
	d.Set("control_plane_id", c.ControlPlaneId)
}

func resourceDPCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newDataPlaneCert := client.DataPlaneCertificate{}
	fillDPCertificate(&newDataPlaneCert, d)
	err := json.NewEncoder(&buf).Encode(newDataPlaneCert)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DataPlanePathCreate, d.Get("control_plane_id").(string))
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.DataPlaneCertificate{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newDataPlaneCert.ControlPlaneId
	d.SetId(retVal.DataPlaneCertificateEncodeId())
	fillResourceDataFromDPCertificate(retVal, d)
	return diags

}

func resourceDPCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.DataPlaneCertificateDecodeId(d.Id())

	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DataPlanePathGet, controlPlaneId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.DataPlaneCertificate{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromDPCertificate(retVal, d)
	return diags
}

func resourceDPCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.CertificateDecodeId(d.Id())

	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DataPlanePathGet, controlPlaneId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
