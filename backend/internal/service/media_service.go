package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/osmanmertacar/elci/backend/internal/config"
)

// MediaService presigns direct-to-R2 uploads: the browser PUTs the file
// straight to object storage, so uploads never stream through this process.
type MediaService struct {
	presign       *s3.PresignClient
	bucket        string
	publicBaseURL string
}

func NewMediaService(ctx context.Context, cfg config.R2Config) (*MediaService, error) {
	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
		o.UsePathStyle = true
	})

	return &MediaService{
		presign:       s3.NewPresignClient(client),
		bucket:        cfg.Bucket,
		publicBaseURL: strings.TrimRight(cfg.PublicBaseURL, "/"),
	}, nil
}

type PresignedUpload struct {
	Key       string
	UploadURL string
	PublicURL string
}

// PresignUpload issues a time-limited PUT URL for a new object under the
// given user's namespace, keyed by a random id so uploads can never collide
// or be guessed.
func (s *MediaService) PresignUpload(ctx context.Context, userID int64, filename, contentType string) (PresignedUpload, error) {
	key := fmt.Sprintf("%d/%s%s", userID, uuid.NewString(), filepath.Ext(filename))

	req, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return PresignedUpload{}, fmt.Errorf("presign upload: %w", err)
	}

	return PresignedUpload{
		Key:       key,
		UploadURL: req.URL,
		PublicURL: s.publicBaseURL + "/" + key,
	}, nil
}
