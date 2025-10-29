package api

import "io"

func (a *APIClient) UploadTSV(f io.ReadSeeker) error {
	return easyNilFetch(a, `POST`, `/upload`, f)
}
