// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceSetNested{}
var _ resource.ResourceWithImportState = &ResourceSetNested{}

func NewResourceSetNested() resource.Resource {
	return &ResourceSetNested{}
}

// ResourceSetNested defines the resource implementation.
type ResourceSetNested struct {
	client *http.Client
}

// ResourceSetNestedModel describes the resource data model.
type (
	ResourceSetNestedModel struct {
		Id        types.String `tfsdk:"id"`
		SetNested types.Set    `tfsdk:"set_nested"`
	}

	SetNestedModel struct {
		Uuid          types.String `tfsdk:"uuid"`
		FixedIp       types.String `tfsdk:"fixed_ip"`
		FixedIpV4     types.String `tfsdk:"fixed_ip_v4"`
		Port          types.String `tfsdk:"port"`
		Mac           types.String `tfsdk:"mac"`
		EnableGateway types.Bool   `tfsdk:"enable_gateway"`
	}
)

var (
	setNestedModelTypeMap = map[string]attr.Type{
		"uuid":           types.StringType,
		"fixed_ip":       types.StringType,
		"fixed_ip_v4":    types.StringType,
		"port":           types.StringType,
		"mac":            types.StringType,
		"enable_gateway": types.BoolType,
	}
)

func (r *ResourceSetNested) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set_nested"
}

func (r *ResourceSetNested) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"set_nested": schema.SetNestedAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							MarkdownDescription: "经典网络ID",
							Required:            true,
						},
						"fixed_ip": schema.StringAttribute{
							MarkdownDescription: "指定IP地址",
							Optional:            true,
							Computed:            true,
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

							//Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *ResourceSetNested) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceSetNested) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ResourceSetNestedModel

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

func (r *ResourceSetNested) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ResourceSetNestedModel

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

func (r *ResourceSetNested) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceSetNestedModel

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

func (r *ResourceSetNested) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceSetNestedModel

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

func (r *ResourceSetNested) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (s *ResourceSetNestedModel) fnConvert(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	sSetNested := s.SetNested
	var sSetNestedModels []SetNestedModel
	sSetNested.ElementsAs(ctx, &sSetNestedModels, false)

	processedSetNestedModels := make([]SetNestedModel, 0)

	for i, model := range sSetNestedModels {
		model.Port = types.StringValue(fmt.Sprintf("port_id_%d", i))
		model.Mac = types.StringValue(fmt.Sprintf("mac_address_%d", i))

		mFixedIp := model.FixedIp
		if mFixedIp.IsNull() || mFixedIp.IsUnknown() {
			model.FixedIp = types.StringValue(fmt.Sprintf("fixed_ip_%d", i))
			model.FixedIpV4 = types.StringValue(fmt.Sprintf("fixed_ip_v4_%d", i))
		} else {
			model.FixedIpV4 = mFixedIp
		}

		processedSetNestedModels = append(processedSetNestedModels, model)
	}

	if len(processedSetNestedModels) > 0 {
		sets, diags1 := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: setNestedModelTypeMap,
		}, processedSetNestedModels)

		diags.Append(diags1...)

		if diags.HasError() {
			return diags
		}

		s.SetNested = sets
	}

	return diags
}
