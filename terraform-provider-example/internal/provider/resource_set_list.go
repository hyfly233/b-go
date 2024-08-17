// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceSetList{}

func NewResourceSetListList() resource.Resource {
	return &ResourceSetList{}
}

// ResourceSetList defines the resource implementation.
type ResourceSetList struct {
	client *http.Client
}

// ResourceSetListModel describes the resource data model.
type (
	ResourceSetListModel struct {
		Id       types.String `tfsdk:"id"`
		TestSet  types.Set    `tfsdk:"test_set"`
		TestList types.List   `tfsdk:"test_list"`
	}
)

func (r *ResourceSetList) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set_list"
}

func (r *ResourceSetList) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Set List Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"test_set": schema.SetAttribute{
				MarkdownDescription: "test set",
				Required:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeBetween(0, 16),
					setvalidator.ValueStringsAre(
						stringvalidator.All(
							stringvalidator.LengthAtMost(255),
						),
					),
				},
			},
			"test_list": schema.ListAttribute{
				MarkdownDescription: "test set",
				Required:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(0, 16),
					listvalidator.ValueStringsAre(
						stringvalidator.All(
							stringvalidator.LengthAtMost(255),
						),
					),
				},
			},
		},
	}
}

func (r *ResourceSetList) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceSetList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ResourceSetListModel

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

func (r *ResourceSetList) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ResourceSetListModel

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

func (r *ResourceSetList) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceSetListModel

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

func (r *ResourceSetList) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceSetListModel

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

func (s *ResourceSetListModel) fnConvert(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	// set ---------------
	stSet := s.TestSet
	setElements := make([]types.String, 0, len(stSet.Elements()))
	diags = stSet.ElementsAs(ctx, &setElements, false)

	if diags.HasError() {
		return diags
	}

	for _, tsElement := range setElements {
		log.Printf("before set element: %s", tsElement.ValueString())

		if !tsElement.IsNull() && !tsElement.IsUnknown() {
			log.Printf("set element: %s", tsElement.ValueString())

			if "abc" == tsElement.ValueString() {
				log.Printf("set element == abc")
			} else {
				log.Printf("set element != abc")
			}
		}

		log.Printf("after set element: %s", tsElement.ValueString())
	}

	// list ---------------
	stList := s.TestList

	listElements := make([]types.String, 0, len(stList.Elements()))
	diags = stList.ElementsAs(ctx, &listElements, false)

	if diags.HasError() {
		return diags
	}

	for _, tlElement := range listElements {
		if !tlElement.IsNull() && !tlElement.IsUnknown() {
			log.Printf("list element: %s", tlElement.ValueString())
		}
	}

	return diags
}
