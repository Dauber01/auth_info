package document

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"

	"auth_info/internal/pkg/apperr"
)

func TestFetchImageBytes_Base64TooLarge(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte("a"), maxImageBytes+1))
	imageURL := "data:image/png;base64," + encoded

	uc := NewUseCase(nil)
	_, err := uc.fetchImageBytes(context.Background(), ImageValue{ImageURL: imageURL})
	if !apperr.IsCode(err, apperr.CodeInvalidArgument) {
		t.Fatalf("expected invalid argument error, got: %v", err)
	}
}
