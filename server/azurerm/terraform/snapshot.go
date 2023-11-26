package terraform

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type snapshotValues struct {
	Location            string    `mapstructure:"location"`
	DiskSizeGB          *float64  `mapstructure:"disk_size_gb"`
	SourceUri           []float64 `mapstructure:"source_uri"`
	DiskSizeGBSourceUri float64   `mapstructure:"disk_size_gb_source_uri"` //this disk sizeGB is inside the SourceUri
}

func decodeSnapshotValues(tfVals map[string]interface{}) (snapshotValues, error) {
	var v snapshotValues
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &v,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return v, err
	}

	if err := decoder.Decode(tfVals); err != nil {
		return v, err
	}
	return v, nil
}

func (p *Provider) newSnapshot(vals snapshotValues) *Image {
	return &Image{
		imageType: "Snapshot",
		storageGB: decimal.NewFromFloat(*snapshotStorageSize(vals)),
		location:  vals.Location,
	}
}

func snapshotStorageSize(vals snapshotValues) *float64 {
	if vals.DiskSizeGB != nil {
		return vals.DiskSizeGB
	}

	if len(vals.SourceUri) > 0 {
		size := vals.DiskSizeGBSourceUri
		return &size
	}

	return nil
}
