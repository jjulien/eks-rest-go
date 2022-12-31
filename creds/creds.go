package creds

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const URL_EXPIRES = 60
const emptySha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func BearerToken(ctx context.Context, clusterName string, cfg aws.Config) (string, error) {
	signer := v4.NewSigner()
	q := url.Values{}
	q.Add("Action", "GetCallerIdentity")
	q.Add("Version", "2011-06-15")
	q.Add("X-Amz-Expires", strconv.FormatInt(int64(URL_EXPIRES/time.Second), 10))
	httpRequest, err := http.NewRequest(http.MethodGet, "https://sts.amazonaws.com/", nil)
	if err != nil {
		return "", err
	}

	httpRequest.URL.RawQuery = q.Encode()
	httpRequest.Header.Set("x-k8s-aws-id", clusterName)

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return "", err
	}

	uri, _, err := signer.PresignHTTP(ctx, creds,
		httpRequest, emptySha256, "sts", cfg.Region, time.Now(), func(*v4.SignerOptions) {})
	if err != nil {
		return "", err
	}

	b64PreSigned := base64.URLEncoding.EncodeToString([]byte(uri))
	eksBearerToken := fmt.Sprintf("k8s-aws-v1.%s", strings.TrimRight(b64PreSigned, "="))

	return eksBearerToken, nil
}
