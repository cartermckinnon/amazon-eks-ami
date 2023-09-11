package configprovider

import (
	"strings"
	"testing"

	internalapi "github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	"github.com/stretchr/testify/assert"
)

func Test_parseMIMEMessage(t *testing.T) {
	testcases := []struct {
		name          string
		input         string
		expected      *internalapi.NodeConfig
		errorExpected bool
	}{
		{
			name: "valid MIME message",
			input: `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="BOUNDARY"

--BOUNDARY
Content-Type: application/json; charset="utf-8"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="nodeconfig.json"

eyJjb25maWciOiAiY29uZmlnIiwgImRhdGEiOiB7Im5hbWUiOiAibm9kZWNvbmZpZyIsICJjb250ZW50IjogIm5vZGVjb25maWcifX0=

--BOUNDARY
Content-Type: node.eks.aws/v1alpha1; charset="utf-8"

---
apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig

--BOUNDARY--`,
			expected: &internalapi.NodeConfig{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseMIMEMultiPart(strings.NewReader(tc.input))
			if err != nil {
				if !tc.errorExpected {
					t.Fatal(err)
				}
			} else {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
