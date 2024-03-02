package bundle

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hamba/avro"
	"github.com/liteseed/bungo/internal/types"
)

func getSignatureMetadata(data []byte, N int) (SignatureType int, SignatureLength int, PublicKeyLength int, err error) {
	SignatureType = int(binary.LittleEndian.Uint16(data))
	signatureMeta, ok := types.SignatureConfig[SignatureType]
	if !ok {
		return -1, -1, -1, fmt.Errorf("unsupported signature type:%d", SignatureType)
	}
	if N < signatureMeta.SignatureLength+2 {
		return -1, -1, -1, errors.New("dataItem longer than expected signature length")
	}
	SignatureLength = signatureMeta.SignatureLength
	PublicKeyLength = signatureMeta.PublicKeyLength
	err = nil
	return
}

func getTarget(data *[]byte, startAt int) (string, int) {
	target := ""
	position := startAt
	if (*data)[startAt] == 1 {
		target = base64.URLEncoding.EncodeToString((*data)[startAt+1 : startAt+1+32])
		position += 32
	}
	return target, position
}

func getAnchor(data *[]byte, startAt int) (string, int) {
	anchor := ""
	position := startAt
	if (*data)[startAt] == 1 {
		anchor = base64.URLEncoding.EncodeToString((*data)[position+1 : position+1+32])
		position += 32
	}
	return anchor, position
}

func decodeTags(data *[]byte, startAt int) (*[]types.Tag, int, error) {
	tags := &[]types.Tag{}
	tagsEnd := startAt + 8

	numberOfTags := int(binary.LittleEndian.Uint16((*data)[startAt : startAt+8]))

	if numberOfTags > 0 {

		numberOfTagBytesStart := startAt + 8
		numberOfTagBytesEnd := numberOfTagBytesStart + 8
		numberOfTagBytes := int(binary.LittleEndian.Uint16((*data)[numberOfTagBytesStart:numberOfTagBytesEnd]))

		println(numberOfTagBytes)
		bytesDataStart := numberOfTagBytesEnd
		bytesDataEnd := numberOfTagBytesEnd + numberOfTagBytes
		bytesData := (*data)[bytesDataStart:bytesDataEnd]

		tags, err := decodeAvro(bytesData)
		if err != nil {
			return nil, tagsEnd, err
		}
		println(tags)
		tagsEnd = bytesDataEnd
		return tags, tagsEnd, nil
	}
	return tags, tagsEnd, nil
}

type avroTag struct {
	name  []byte
	value []byte
}

func decodeAvro(data []byte) (*[]types.Tag, error) {
	codec, err := avro.Parse(`
	{
		"type": "array",
		"items": {
			"type": "record",
			"name": "Tag",
			"fields": [
				{ "name": "name", "type": "bytes" },
				{ "name": "value", "type": "bytes" }
			]
		}
	}`)
	if err != nil {
		panic(err)
	}
	avroTags := &[]map[string]any{}
	err = avro.Unmarshal(codec, data, avroTags)
	if err != nil {
		return nil, err
	}

	tags := []types.Tag{}
	for _, v := range *avroTags {
		tags = append(tags, types.Tag{Name: string(v["name"].([]byte)), Value: string(v["value"].([]byte))})
	}
	return &tags, err
}

// Decode a DataItem from bytes
func DecodeDataItem(data []byte) (*types.DataItem, error) {
	N := len(data)
	if N < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(data[:2], N)
	if err != nil {
		return nil, err
	}
	signatureStart := 2
	signatureEnd := signatureLength + signatureStart
	signature := base64.URLEncoding.EncodeToString(data[signatureStart:signatureEnd])

	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := base64.URLEncoding.EncodeToString(data[ownerStart:ownerEnd])

	position := 2 + ownerEnd
	target, position := getTarget(&data, position)
	anchor, position := getAnchor(&data, position)
	tags, position, err := decodeTags(&data, position)
	if err != nil {
		return nil, err
	}

	rawData := data[position:]

	return &types.DataItem{
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		RawData:       base64.URLEncoding.EncodeToString(rawData),
	}, nil
}

func DecodeBundle(data []byte) (*types.Bundle, error) {
	// length must more than 32
	if len(data) < 32 {
		return nil, errors.New("binary length must more than 32")
	}
	N := int(binary.LittleEndian.Uint16(data[:32]))

	if len(data) < 32+N*64 {
		return nil, errors.New("binary length incorrect")
	}

	bundle := &types.Bundle{
		Items:   make([]types.DataItem, 0),
		RawData: data,
	}
	bundleItemStart := 32 + N*64
	for i := 0; i < N; i++ {
		headerBegin := 32 + i*64
		end := headerBegin + 64
		if len(data) < end {
			return nil, errors.New("binary length incorrect")
		}
		headerByte := data[headerBegin:end]
		itemBinaryLength := int(binary.LittleEndian.Uint16(headerByte[:32]))
		id := base64.URLEncoding.EncodeToString(headerByte[32:64])
		if len(data) < bundleItemStart+itemBinaryLength || itemBinaryLength < 0 {
			return nil, errors.New("binary length incorrect")
		}
		bundleItemBytes := data[bundleItemStart : bundleItemStart+itemBinaryLength]
		bundleItem, err := DecodeDataItem(bundleItemBytes)
		if err != nil {
			return nil, err
		}
		if bundleItem.Id != id {
			return nil, fmt.Errorf("bundleItem.Id != id, bundleItem.Id: %s, id: %s", bundleItem.Id, id)
		}
		bundle.Items = append(bundle.Items, *bundleItem)
		bundleItemStart += itemBinaryLength
	}
	return bundle, nil
}
