package av1

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.AV1_OBU,
		Description: "AV1 Open Bitstream Unit",
		DecodeFn:    obuDecode,
	})
}

//nolint:revive
const (
	OBU_SEQUENCE_HEADER        = 1
	OBU_TEMPORAL_DELIMITER     = 2
	OBU_FRAME_HEADER           = 3
	OBU_TILE_GROUP             = 4
	OBU_METADATA               = 5
	OBU_FRAME                  = 6
	OBU_REDUNDANT_FRAME_HEADER = 7
	OBU_TILE_LIST              = 8
	OBU_PADDING                = 15
)

var obuTypeNames = scalar.UintMapSymStr{
	OBU_SEQUENCE_HEADER:        "OBU_SEQUENCE_HEADER",
	OBU_TEMPORAL_DELIMITER:     "OBU_TEMPORAL_DELIMITER",
	OBU_FRAME_HEADER:           "OBU_FRAME_HEADER",
	OBU_TILE_GROUP:             "OBU_TILE_GROUP",
	OBU_METADATA:               "OBU_METADATA",
	OBU_FRAME:                  "OBU_FRAME",
	OBU_REDUNDANT_FRAME_HEADER: "OBU_REDUNDANT_FRAME_HEADER",
	OBU_TILE_LIST:              "OBU_TILE_LIST",
	OBU_PADDING:                "OBU_PADDING",
}

func obuDecode(d *decode.D) any {
	var obuType uint64
	var obuSize int64
	hasExtension := false
	hasSizeField := false

	d.FieldStruct("header", func(d *decode.D) {
		d.FieldU1("forbidden_bit")
		obuType = d.FieldU4("type", obuTypeNames)
		hasExtension = d.FieldBool("extension_flag")
		hasSizeField = d.FieldBool("has_size_field")
		d.FieldU1("reserved_1bit")
		if hasExtension {
			d.FieldU3("temporal_id")
			d.FieldU2("spatial_id")
			d.FieldU3("extension_header_reserved_3bits")
		}
	})

	if hasSizeField {
		obuSize = int64(d.FieldULEB128("size"))
	} else {
		obuSize = d.BitsLeft() / 8
		if hasExtension {
			obuSize--
		}
	}

	_ = obuType

	if d.BitsLeft() > 0 {
		d.FieldRawLen("data", obuSize*8)
	}

	return nil
}
