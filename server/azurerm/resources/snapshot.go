package resources

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type SourceUri struct {
	Values struct {
		DiskSizeGb *float64 `mapstructure:"disk_size_gb"`
	} `mapstructure:"values"`
}

type snapshotValues struct {
	Location   string      `mapstructure:"location"`
	DiskSizeGB *float64    `mapstructure:"disk_size_gb"`
	SourceUri  []SourceUri `mapstructure:"source_uri"`
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
	var storageGB float64
	storageSize := snapshotStorageSize(vals)
	if storageSize != nil {
		storageGB = *storageSize
	}
	return &Image{
		imageType: "Snapshot",
		storageGB: decimal.NewFromFloat(storageGB),
		location:  getLocationName(vals.Location),
	}
}

func snapshotStorageSize(vals snapshotValues) *float64 {
	if vals.DiskSizeGB != nil {
		return vals.DiskSizeGB
	}

	if vals.SourceUri != nil {
		size := vals.SourceUri[0].Values.DiskSizeGb
		return size
	}

	return nil
}
