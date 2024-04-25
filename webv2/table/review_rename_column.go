// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package table

import (
	"github.com/GoogleCloudPlatform/spanner-migration-tool/internal"
	utilities "github.com/GoogleCloudPlatform/spanner-migration-tool/webv2/utilities"
)

// reviewRenameColumn review  renaming of Columnname in schmema.
func reviewRenameColumn(newName, tableId, colId string, conv *internal.Conv, interleaveTableSchema []InterleaveTableSchema) []InterleaveTableSchema {

	sp := conv.SpSchema.Tables[tableId]

	// review column name update for interleaved child.
	interleaveTableSchema, childTableId := reviewRenameColumnForChildTable(newName, tableId, colId, conv, interleaveTableSchema)

	// review column name update for interleaved parent.
	interleaveTableSchema, parentTableId := reviewRenameColumnForParentTable(newName, tableId, colId, conv, interleaveTableSchema)

	oldColName := conv.SpSchema.Tables[tableId].ColDefs[colId].Name
	colType := conv.SpSchema.Tables[tableId].ColDefs[colId].T.Name
	colSize := conv.SpSchema.Tables[tableId].ColDefs[colId].T.Len
	reviewRenameColumnNameTableSchema(conv, tableId, colId, newName)
	if childTableId != "" || parentTableId != "" {
		interleaveTableSchema = renameInterleaveTableSchema(interleaveTableSchema, sp.Name, colId, oldColName, newName, colType, int(colSize))
	}

	return interleaveTableSchema
}

func reviewRenameColumnForChildTable(newName, tableId, colId string, conv *internal.Conv, interleaveTableSchema []InterleaveTableSchema) ([]InterleaveTableSchema, string) {
	sp := conv.SpSchema.Tables[tableId]
	isParent, childTableId := utilities.IsParent(tableId)

	if isParent {
		childColId, err := utilities.GetColIdFromSpannerName(conv, childTableId, sp.ColDefs[colId].Name)
		if err == nil {
			interleaveTableSchema, _ = reviewRenameColumnForChildTable(newName, childTableId, childColId, conv, interleaveTableSchema)
			oldColName := conv.SpSchema.Tables[childTableId].ColDefs[childColId].Name
			reviewRenameColumnNameTableSchema(conv, childTableId, childColId, newName)
			childTableName := conv.SpSchema.Tables[childTableId].Name
			colType := conv.SpSchema.Tables[childTableId].ColDefs[childColId].T.Name
			colSize := conv.SpSchema.Tables[childTableId].ColDefs[childColId].T.Len
			interleaveTableSchema = renameInterleaveTableSchema(interleaveTableSchema, childTableName, childColId, oldColName, newName, colType, int(colSize))
		}
	}
	return interleaveTableSchema, childTableId
}

func reviewRenameColumnForParentTable(newName, tableId, colId string, conv *internal.Conv, interleaveTableSchema []InterleaveTableSchema) ([]InterleaveTableSchema, string) {
	sp := conv.SpSchema.Tables[tableId]
	parentTableId := conv.SpSchema.Tables[tableId].ParentId

	if parentTableId != "" {
		parentColId, err := utilities.GetColIdFromSpannerName(conv, parentTableId, sp.ColDefs[colId].Name)
		if err == nil {
			interleaveTableSchema, _ = reviewRenameColumnForParentTable(newName, parentTableId, parentColId, conv, interleaveTableSchema)
			oldColName := conv.SpSchema.Tables[parentTableId].ColDefs[parentColId].Name
			reviewRenameColumnNameTableSchema(conv, parentTableId, parentColId, newName)
			parentTableName := conv.SpSchema.Tables[parentTableId].Name
			colType := conv.SpSchema.Tables[parentTableId].ColDefs[parentColId].T.Name
			colSize := conv.SpSchema.Tables[parentTableId].ColDefs[parentColId].T.Len
			interleaveTableSchema = renameInterleaveTableSchema(interleaveTableSchema, parentTableName, parentColId, oldColName, newName, colType, int(colSize))
		}
	}
	return interleaveTableSchema, parentTableId
}

// reviewRenameColumnNameTableSchema review  renaming of column-name in Table Schema.
func reviewRenameColumnNameTableSchema(conv *internal.Conv, tableId, colId, newName string) {
	sp := conv.SpSchema.Tables[tableId]

	column, ok := sp.ColDefs[colId]

	if ok {
		column.Name = newName

		sp.ColDefs[colId] = column
		conv.SpSchema.Tables[tableId] = sp

	}
}
