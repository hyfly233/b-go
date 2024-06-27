// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
var _ resource.Resource = &RegexResource{}
var _ resource.ResourceWithImportState = &RegexResource{}

func NewRegexResource() resource.Resource {
	return &RegexResource{}
}

// RegexResource defines the resource implementation.
type RegexResource struct {
	client *http.Client
}

// RegexResourceModel describes the resource data model.
type (
	RegexResourceModel struct {
		Id    types.String `tfsdk:"id"`
		Name  types.String `tfsdk:"name"`
		Alias types.String `tfsdk:"alias"`
	}
)

func (r *RegexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regex"
}

func (r *RegexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"name": schema.StringAttribute{
				MarkdownDescription: "虚机名称",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(25),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$`),
						"只能使用字母、数字和短横线，且必须以字母开头，不能以短横线结尾",
					),
				},
			},
			"alias": schema.StringAttribute{
				MarkdownDescription: "虚机别名",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(0),
					stringvalidator.LengthAtMost(32),
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[^.\u3000-\u303F\\/:*?\"<>|][^\\/:*?\"<>|]{0,32}[^.\u3000-\u303F\\/:*?\"<>|]?$"),
						"有效长度为0至32个字符，不能包含全角字符以及括号中的英文字符（/：*?\"<>|），并且不能以点'.'作为开始和结束字符",
					),
				},
			},
		},
	}
}

func (r *RegexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RegexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RegexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RegexResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RegexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RegexResourceModel

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

func (r *RegexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
