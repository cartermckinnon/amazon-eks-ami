package configprovider

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/awslabs/amazon-eks-ami/nodeadm/api"
	internalapi "github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	apibridge "github.com/awslabs/amazon-eks-ami/nodeadm/internal/api/bridge"
)

type imdsConfigProvider struct {
	client *imds.Client
}

func NewIMDSConfigProvider(path string) ConfigProvider {
	return &imdsConfigProvider{
		client: imds.New(imds.Options{}),
	}
}

func (p *imdsConfigProvider) Provide() (*internalapi.NodeConfig, error) {
	resp, err := p.client.GetUserData(context.TODO(), &imds.GetUserDataInput{})
	if err != nil {
		return nil, err
	}
	return parseMIMEMultiPart(resp.Content)
}

func parseMIMEMultiPart(data io.Reader) (*internalapi.NodeConfig, error) {
	msg, err := mail.ReadMessage(data)
	if err != nil {
		return nil, err
	}
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, fmt.Errorf("invalid MIME media type: %s", mediaType)
	}
	multiPartReader := multipart.NewReader(msg.Body, params["boundary"])
	for {
		part, err := multiPartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		nodeConfig, err := parseMIMEPart(part)
		if err == nil {
			return nodeConfig, nil
		}
	}
	return nil, fmt.Errorf("no MIME part with %s media type found", api.GroupName)
}

func parseMIMEPart(part *multipart.Part) (*internalapi.NodeConfig, error) {
	mediaType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(mediaType, api.GroupName+"/") {
		content, err := decodeMIMEPart(part)
		if err != nil {
			return nil, err
		}
		return apibridge.DecodeNodeConfig(content)
	}
	return nil, fmt.Errorf("unknown MIME media type: %s", mediaType)
}

func decodeMIMEPart(part *multipart.Part) ([]byte, error) {
	content, err := io.ReadAll(part)
	if err != nil {
		return nil, err
	}
	contentTransferEncoding := part.Header.Get("Content-Transfer-Encoding")
	switch contentTransferEncoding {
	case "base64":
		decodedContent, err := base64.StdEncoding.DecodeString(string(content))
		if err != nil {
			return nil, err
		}
		return decodedContent, nil
	case "quoted-printable":
		decoded, err := io.ReadAll(quotedprintable.NewReader(bytes.NewReader(content)))
		if err != nil {
			return nil, err
		}
		return decoded, nil
	default:
		return content, nil
	}
}
