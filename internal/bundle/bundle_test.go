package bundle

import (
	"log"
	"os"
	"testing"

	"github.com/liteseed/bungo/internal/types"
	"gotest.tools/v3/assert"
)

func TestDecodeDataItem(t *testing.T) {
	data, err := os.ReadFile("../../test/stubs/1115BDataItem")
	if err != nil {
		log.Fatal(err)
	}
	dataItem, err := DecodeDataItem(data)
	assert.NilError(t, err)
	assert.Equal(t, dataItem.Signature, "wUIlPaBflf54QyfiCkLnQcfakgcS5B4Pld-hlOJKyALY82xpAivoc0fxBJWjoeg3zy9aXz8WwCs_0t0MaepMBz2bQljRrVXnsyWUN-CYYfKv0RRglOl-kCmTiy45Ox13LPMATeJADFqkBoQKnGhyyxW81YfuPnVlogFWSz1XHQgHxrFMAeTe9epvBK8OCnYqDjch4pwyYUFrk48JFjHM3-I2kcQnm2dAFzFTfO-nnkdQ7ulP3eoAUr-W-KAGtPfWdJKFFgWFCkr_FuNyHYQScQo-FVOwIsvj_PVWEU179NwiqfkZtnN8VoBgCSxbL1Wmh4NYL-GsRbKz_94hpcj5RiIgq0_H5dzAp-bIb49M4SP-DcuIJ5oT2v2AfPWvznokDDVTeikQJxCD2n9usBOJRpLw_P724Yurbl30eNow0U-Jmrl8S6N64cjwKVLI-hBUfcpviksKEF5_I4XCyciW0TvZj1GxK6ET9lx0s6jFMBf27-GrFx6ZDJUBncX6w8nDvuL6A8TG_ILGNQU_EDoW7iil6NcHn5w11yS_yLkqG6dw_zuC1Vkg1tbcKY3703tmbF-jMEZUvJ6oN8vRwwodinJjzGdj7bxmkUPThwVWedCc8wCR3Ak4OkIGASLMUahSiOkYmELbmwq5II-1Txp2gDPjCpAf9gT6Iu0heAaXhjk=")
	assert.Equal(t, dataItem.Owner, "0zBGbs8Y4wvdS58cAVyxp7mDffScOkbjh50ZrqnWKR_5NGwjezT6J40ejIg5cm1KnuDnw9OhvA7zO6sv1hEE6IaGNnNJWiXFecRMxCl7iw78frrT8xJvhBgtD4fBCV7eIvydqLoMl8K47sacTUxEGseaLfUdYVJ5CSock5SktEEdqqoe3MAso7x4ZsB5CGrbumNcCTifr2mMsrBytocSoHuiCEi7-Nwv4CqzB6oqymBtEECmKYWdINnNQHVyKK1l0XP1hzByHv_WmhouTPos9Y77sgewZrvLF-dGPNWSc6LaYGy5IphCnq9ACFrEbwkiCRgZHnKsRFH0dfGaCgGb3GZE-uspmICJokJ9CwDPDJoxkCBEF0tcLSIA9_ofiJXaZXbrZzu3TUXWU3LQiTqYr4j5gj_7uTclewbyZSsY-msfbFQlaACc02nQkEkr4pMdpEOdAXjWP6qu7AJqoBPNtDPBqWbdfsLXgyK90NbYmf3x4giAmk8L9REy7SGYugG4VyqG39pNQy_hdpXdcfyE0ftCr5tSHVpMreJ0ni7v3IDCbjZFcvcHp0H6f6WPfNCoHg1BM6rHUqkXWd84gdHUzo9LTGq9-7wSBCizpcc_12_I-6yvZsROJvdfYOmjPnd5llefa_X3X1dVm5FPYFIabydGlh1Vs656rRu4dzeEQwc=")
	assert.Equal(t, dataItem.Target, "")
	assert.Equal(t, dataItem.Anchor, "")
	assert.DeepEqual(
		t,
		dataItem.Tags,
		[]types.Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "App-Name", Value: "ArDrive-CLI"},
			{Name: "App-Version", Value: "1.21.0"},
		},
	)
	assert.Equal(t, dataItem.RawData, "NTY3MAo=")
}

func TestDecodeBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../../test/stubs/bundleHeader")
	if err != nil {
		log.Fatal(err)
	}
	headers, N := decodeBundleHeader(&data)
	assert.Equal(t, N, 1)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}
