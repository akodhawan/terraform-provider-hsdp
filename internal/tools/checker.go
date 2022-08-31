package tools

import (
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/philips-software/go-hsdp-api/iam"
)

func CheckForPermissionErrors(client iam.TokenRefresher, status iam.HTTPStatus, err error) error {
	if status != nil && status.StatusCode() > 500 {
		return err
	}
	if status != nil && status.StatusCode() == http.StatusForbidden {
		_ = client.TokenRefresh()
		return err
	}
	return backoff.Permanent(err)
}
