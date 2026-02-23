// SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"terraform/terraform-provider/provider/common/attribute"
	schemabuilder "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type FileSchema struct{}

var _ schemabuilder.ResourceSchema = &FileSchema{}

func (f *FileSchema) GetSchema() schema.Schema {

	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("A datafile that is automatically copied/distributed to defined Resources.")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name/symbol for the object within the platform and the op language (must be unique, only alphanumeric/underscore).",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("destination_path", schema.StringAttribute{
		MarkdownDescription: "Target location for a copied distributed File object.",
		Required:            true,
	})

	builder.AddAttribute("resource_query", schema.StringAttribute{
		MarkdownDescription: "A set of Resources (e.g. host, pod, container), optionally filtered on tags or dynamic conditions.",
		Required:            true,
	})

	// Optional attributes
	builder.AddAttribute("description", schema.StringAttribute{
		MarkdownDescription: "A user-friendly explanation of an object.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "If the object is currently enabled or disabled.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	// File source attributes (mutually exclusive)
	builder.AddAttribute("input_file", schema.StringAttribute{
		MarkdownDescription: "The local source of a distributed File object. (conflicts with inline_data)",
		Optional:            true,
		Computed:            true,
		Validators: []validator.String{
			// Validate only this attribute or inline_data is configured.
			//
			// Note: For the "input_file" to be in the error message, path.MatchRoot("input_file") is included
			// as an argument to the ExactlyOneOf function (although it is not necessary for the validation to work)
			stringvalidator.ExactlyOneOf(path.MatchRoot("input_file"), path.MatchRoot("inline_data")),
		},
	})

	builder.AddAttribute("inline_data", schema.StringAttribute{
		MarkdownDescription: "The inline file data of a distributed File object. (conflicts with input_file)",
		Optional:            true,
		Computed:            true,
	})

	builder.AddAttribute("md5", schema.StringAttribute{
		MarkdownDescription: "The md5 checksum of a file, e.g. filemd5(\"${path.module}/data/example-file.txt\"). It's used to trigger the file upload when the file changes.",
		Optional:            true,
		Computed:            true,
	})

	// File system attributes with version constraints
	builder.AddAttribute("mode", schema.StringAttribute{
		MarkdownDescription: "The File's permissions, like 'chmod', in octal (e.g. '0644').",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("owner", schema.StringAttribute{
		MarkdownDescription: "The File's ownership, like 'chown' (e.g. 'user:group').",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	// Computed attributes
	builder.AddAttribute("file_data", schema.StringAttribute{
		MarkdownDescription: "The destination URL where the file was uploaded to.",
		Computed:            true,
	})

	builder.AddAttribute("file_length", schema.Int64Attribute{
		MarkdownDescription: "Length, in bytes, of a distributed File object.",
		Computed:            true,
	})

	builder.AddAttribute("checksum", schema.StringAttribute{
		MarkdownDescription: "Cryptographic hash (e.g. md5) of a File Resource. This is used to verify the integrity of the file when it is uploaded.",
		Computed:            true,
	})

	return builder.Build()
}

func (f *FileSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"mode": {
			MinVersion: "release-23.0.0",
		},
		"owner": {
			MinVersion: "release-23.0.0",
		},
	}
}
