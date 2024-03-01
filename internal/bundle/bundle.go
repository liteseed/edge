package bundle

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/liteseed/bungo/internal/types"
)

const (
	ARWEAVE_SIGNATURE_LENGTH = 512
)

func getSignatureMetadata(data *[]byte) (SignatureType int, SignatureLength int, PublicKeyLength int, err error) {
	SignatureType = int(binary.BigEndian.Uint16((*data)[:2]))
	signatureMeta, ok := types.SignatureConfig[SignatureType]
	if !ok {
		return -1, -1, -1, fmt.Errorf("unsupported signature type:%d", SignatureType)
	}
	if len(*data) < signatureMeta.SignatureLength+2 {
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
		target = base64.StdEncoding.EncodeToString((*data)[startAt+1 : startAt+1+32])
		position += 32
	}
	return target, position
}

func getAnchor(data *[]byte, startAt int) (string, int) {
	anchor := ""
	position := startAt
	if (*data)[startAt] == 1 {
		anchor = base64.StdEncoding.EncodeToString((*data)[position+1 : position+1+32])
		position += 32
	}
	return anchor, position
}

// Decode a DataItem from bytes
func DecodeDataItem(itemBinary []byte) (*types.DataItem, error) {
	if len(itemBinary) < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(&itemBinary)
	if err != nil {
		return nil, err
	}

	signature := base64.StdEncoding.EncodeToString(itemBinary[2 : signatureLength+2])
	owner := base64.StdEncoding.EncodeToString(itemBinary[signatureLength+2 : signatureLength+2+publicKeyLength])

	position := 2 + signatureLength + publicKeyLength
	target, position := getTarget(&itemBinary, position)
	anchor, position := getAnchor(&itemBinary, position)

	tagsStart := position + 2
	numOfTags := binary.BigEndian.Uint16(itemBinary[tagsStart : tagsStart+8])

	var tagsBytesLength uint32
	tags := &[]types.Tag{}
	tagsBytes := make([]byte, 0)
	if numOfTags > 0 {

		tagsBytesLength = binary.BigEndian.Uint32(itemBinary[tagsStart+8 : tagsStart+16])

		tagsBytes = itemBinary[tagsStart+16 : tagsStart+16+int(tagsBytesLength)]
		// parser tags
		err := json.Unmarshal(tagsBytes, tags)
		if err != nil {
			return nil, err
		}
	}

	data := itemBinary[tagsStart+16+int(tagsBytesLength):]

	return &types.DataItem{
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		Data:          base64.StdEncoding.EncodeToString(data),
	}, nil
}
