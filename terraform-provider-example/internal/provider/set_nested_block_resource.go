// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SetNestedBlockResource{}
var _ resource.ResourceWithImportState = &SetNestedBlockResource{}

func NewSetNestedBlockResource() resource.Resource {
	return &SetNestedBlockResource{}
}

// SetNestedBlockResource defines the resource implementation.
type SetNestedBlockResource struct {
	client *http.Client
}

// SetNestedBlockResourceModel describes the resource data model.
type (
	SetNestedBlockResourceModel struct {
		Id        types.String `tfsdk:"id"`
		SetNested types.Set    `tfsdk:"set_nested"`
	}

	SetNestedBlockModel struct {
		Uuid          types.String `tfsdk:"uuid"`
		FixedIp       types.String `tfsdk:"fixed_ip"`
		FixedIpV4     types.String `tfsdk:"fixed_ip_v4"`
		Port          types.String `tfsdk:"port"`
		Mac           types.String `tfsdk:"mac"`
		EnableGateway types.Bool   `tfsdk:"enable_gateway"`
	}
)

var (
	SetNestedBlockModelTypeMap = map[string]attr.Type{
		"uuid":           types.StringType,
		"fixed_ip":       types.StringType,
		"fixed_ip_v4":    types.StringType,
		"port":           types.StringType,
		"mac":            types.StringType,
		"enable_gateway": types.BoolType,
	}
)

func (r *SetNestedBlockResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set_nested_block"
}

func (r *SetNestedBlockResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Set Nested Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},

		Blocks: map[string]schema.Block{
			"set_nested": schema.SetNestedBlock{
				MarkdownDescription: "Example configurable attribute",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							MarkdownDescription: "经典网络ID",
							Required:            true,
						},
						"fixed_ip": schema.StringAttribute{
							MarkdownDescription: "指定IP地址",
							Optional:            true,
						},
						"fixed_ip_v4": schema.StringAttribute{
							MarkdownDescription: "指定IPv4地址",
							Computed:            true,
						},
						"port": schema.StringAttribute{
							MarkdownDescription: "网卡端口ID",
							Computed:            true,
						},
						"mac": schema.StringAttribute{
							MarkdownDescription: "MAC地址",
							Computed:            true,
						},
						"enable_gateway": schema.BoolAttribute{
							MarkdownDescription: "是否启用网关",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func (r *SetNestedBlockResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SetNestedBlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SetNestedBlockResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue("example-id")
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

func (r *SetNestedBlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SetNestedBlockResourceModel

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

func (r *SetNestedBlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SetNestedBlockResourceModel

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

func (r *SetNestedBlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SetNestedBlockResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *SetNestedBlockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (s *SetNestedBlockResourceModel) fnConvert(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	sSetNested := s.SetNested
	var sSetNestedBlockModels []SetNestedBlockModel
	sSetNested.ElementsAs(ctx, &sSetNestedBlockModels, false)

	processedSetNestedBlockModels := make([]SetNestedBlockModel, 0)

	for i, model := range sSetNestedBlockModels {
		model.Port = types.StringValue(fmt.Sprintf("port_id_%d", i))
		model.Mac = types.StringValue(fmt.Sprintf("mac_address_%d", i))

		mFixedIp := model.FixedIp
		if mFixedIp.IsNull() || mFixedIp.IsUnknown() {
			model.FixedIpV4 = types.StringValue(fmt.Sprintf("fixed_ip_%d", i))
		} else {
			model.FixedIpV4 = mFixedIp
		}

		processedSetNestedBlockModels = append(processedSetNestedBlockModels, model)
	}

	if len(processedSetNestedBlockModels) > 0 {
		sets, diags1 := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: SetNestedBlockModelTypeMap,
		}, processedSetNestedBlockModels)

		diags.Append(diags1...)

		s.SetNested = sets
	}

	return diags
}
