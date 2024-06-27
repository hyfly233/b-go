// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ModifierResource{}
var _ resource.ResourceWithImportState = &ModifierResource{}

func NewModifierResource() resource.Resource {
	return &ModifierResource{}
}

// ModifierResource defines the resource implementation.
type ModifierResource struct {
	client *http.Client
}

// ModifierResourceModel describes the resource data model.
type (
	ModifierResourceModel struct {
		Id               types.String `tfsdk:"id"`
		Name             types.String `tfsdk:"name"`
		TenantId         types.String `tfsdk:"tenant_id"`
		Description      types.String `tfsdk:"description"`
		NetworkId        types.String `tfsdk:"external_network_id"`
		EnableSnat       types.Bool   `tfsdk:"external_enable_snat"`
		ExternalFixedIps types.Set    `tfsdk:"external_fixed_ips"`
	}

	ExternalFixedIpsModel struct {
		SubnetId  types.String        `tfsdk:"subnet_id"`
		IpAddress iptypes.IPv4Address `tfsdk:"ip_address"`
	}
)

var (
	externalFixedIpAttrTypes = map[string]attr.Type{
		"subnet_id":  types.StringType,
		"ip_address": iptypes.IPv4AddressType{},
	}
)

func (r *ModifierResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_modifier"
}

func (r *ModifierResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Set Nested 2 Example resource",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "租户ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "路由器名称",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[^.\u3000-\u303F\\/:*?\"<>|][^\\/:*?\"<>|]{1,32}[^.\u3000-\u303F\\/:*?\"<>|]?$"),
						"有效长度为0至32个字符，不能包含全角字符以及括号中的英文字符（/：*?\"<>|），并且不能以点'.'作为开始和结束字符",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "描述",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			// 外部网关配置
			"external_network_id": schema.StringAttribute{
				MarkdownDescription: "弹性网络ID",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_enable_snat": schema.BoolAttribute{
				MarkdownDescription: "是否开启SNAT",
				Optional:            true,
			},
			"external_fixed_ips": schema.SetNestedAttribute{
				MarkdownDescription: "外部网关IP地址",
				Optional:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"subnet_id": schema.StringAttribute{
							MarkdownDescription: "子网ID",
							Optional:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: "IP地址",
							Optional:            true,
							CustomType:          iptypes.IPv4AddressType{},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *ModifierResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ModifierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ModifierResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue("f91f202e-abe3-40b6-9d7d-f35fb3bf0471")
	diags := data.fnConvert(ctx)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ModifierResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ModifierResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := data.fnConvert(ctx)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ModifierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ModifierResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := data.fnConvert(ctx)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ModifierResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ModifierResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ModifierResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (s *ModifierResourceModel) fnConvert(ctx context.Context) diag.Diagnostics {

	detail := &RouterRespStruct{
		Id:          StringPtr("f91f202e-abe3-40b6-9d7d-f35fb3bf0471"),
		TenantId:    s.TenantId.ValueStringPointer(),
		Name:        s.Name.ValueStringPointer(),
		Description: s.Description.ValueStringPointer(),
		Status:      StringPtr("ACTIVE"),

		ExternalGatewayInfo: &ExternalGatewayStruct{
			NetworkId:  s.NetworkId.ValueStringPointer(),
			EnableSnat: BoolPtr(false),
			ExternalFixedIps: []*ExternalFixedIpsReqStruct{
				{
					SubnetId:  StringPtr("a15a53f3-85b2-4e07-83b9-8d7a1a672451"),
					IpAddress: StringPtr("127.0.0.1"),
				},
			},
		},
	}

	s.Id = types.StringPointerValue(detail.Id)
	s.Name = types.StringPointerValue(detail.Name)
	s.TenantId = types.StringPointerValue(detail.TenantId)
	s.Description = types.StringPointerValue(detail.Description)

	// external_gateway_info
	if egInfo := detail.ExternalGatewayInfo; egInfo != nil {
		s.NetworkId = types.StringPointerValue(egInfo.NetworkId)
		s.EnableSnat = types.BoolPointerValue(egInfo.EnableSnat)

		// external_fixed_ips
		sFixedIps := egInfo.ExternalFixedIps

		elements := make([]*ExternalFixedIpsModel, 0)
		for _, ip := range sFixedIps {
			elements = append(elements, &ExternalFixedIpsModel{
				SubnetId:  types.StringPointerValue(ip.SubnetId),
				IpAddress: iptypes.NewIPv4AddressPointerValue(ip.IpAddress),
			})
		}

		sets, diags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: externalFixedIpAttrTypes,
		}, elements)

		if !diags.HasError() {
			s.ExternalFixedIps = sets
		} else {
			return diags
		}
	}

	return nil

}

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

type (
	ExternalGatewayStruct struct {
		NetworkId        *string                      `json:"network_id,omitempty"`
		EnableSnat       *bool                        `json:"enable_snat,omitempty"`
		ExternalFixedIps []*ExternalFixedIpsReqStruct `json:"external_fixed_ips,omitempty"`
	}

	ExternalFixedIpsReqStruct struct {
		SubnetId  *string `json:"subnet_id,omitempty"`
		IpAddress *string `json:"ip_address,omitempty"`
	}

	RouterDetailStruct struct {
		Router *RouterRespStruct `json:"router,omitempty"`
	}

	RouterRespStruct struct {
		Id                    *string                `json:"id,omitempty"`
		TenantId              *string                `json:"tenant_id,omitempty"`
		ProjectId             *string                `json:"project_id,omitempty"`
		Name                  *string                `json:"name,omitempty"`
		Description           *string                `json:"description,omitempty"`
		Status                *string                `json:"status,omitempty"`
		Ha                    *bool                  `json:"ha,omitempty"`
		AdminStateUp          *bool                  `json:"admin_state_up,omitempty"`
		Distributed           *bool                  `json:"distributed,omitempty"`
		ExternalGatewayInfo   *ExternalGatewayStruct `json:"external_gateway_info,omitempty"`
		Routes                []*RoutesItemStruct    `json:"routes,omitempty"`
		AvailabilityZoneHints []interface{}          `json:"availability_zone_hints,omitempty"`
		AvailabilityZones     []interface{}          `json:"availability_zones,omitempty"`
		Tags                  []interface{}          `json:"tags,omitempty"`
		FlavorId              interface{}            `json:"flavor_id,omitempty"`
	}

	RoutesItemStruct struct {
		NextHop     *string `json:"nexthop"`
		Destination *string `json:"destination"`
	}
)
