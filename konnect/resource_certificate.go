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

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		UpdateContext: resourceCertificateUpdate,
		DeleteContext: resourceCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"certificate": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, math.MaxInt),
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alternate_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alternate_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func fillCertificate(c *client.Certificate, d *schema.ResourceData) {
	c.Certificate = d.Get("certificate").(string)
	c.Key = d.Get("key").(string)
	alternateKey, ok := d.GetOk("alternate_key")
	if ok {
		c.AlternateKey = alternateKey.(string)
	}
	alternateCertificate, ok := d.GetOk("alternate_certificate")
	if ok {
		c.AlternateCertificate = alternateCertificate.(string)
	}
}

func fillResourceDataFromCertificate(c *client.Certificate, d *schema.ResourceData) {
	d.Set("certificate", c.Certificate)
	d.Set("key", c.Key)
	d.Set("alternate_key", c.AlternateKey)
	d.Set("alternate_certificate", c.AlternateCertificate)
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newDataPlaneCert := client.Certificate{}
	fillCertificate(&newDataPlaneCert, d)
	err := json.NewEncoder(&buf).Encode(newDataPlaneCert)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CertificatePathCreate, d.Get("control_plane_id").(string))
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Certificate{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Id)
	fillResourceDataFromCertificate(retVal, d)
	d.Set("control_plane_id", d.Get("control_plane_id").(string))
	return diags
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CertificatePathGet, d.Get("control_plane_id").(string), d.Id())
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Certificate{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromCertificate(retVal, d)
	return diags
}

func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CertificatePathGet, d.Get("control_plane_id").(string), d.Id())
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCertificate := client.Certificate{}
	fillCertificate(&upCertificate, d)
	err := json.NewEncoder(&buf).Encode(upCertificate)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CertificatePathGet, d.Get("control_plane_id").(string), d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Certificate{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceDataFromCertificate(retVal, d)
	return diags
}
