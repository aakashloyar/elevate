package s3

import (
	"context"
	"io"
	"time"

	"github.com/aakashloyar/elevate/document/config"
	"github.com/aakashloyar/elevate/document/internal/domain"
	"github.com/aakashloyar/elevate/document/internal/observe"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucketName string
}

func NewS3Storage(client *s3.Client, bucket string) *S3Storage {
	return &S3Storage{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		bucketName: bucket,
	}
}

func (s *S3Storage) BuildStorageKey(uploadID, fileName string) string {
	return "uploads/" + uploadID + "/" + fileName
}

func (s *S3Storage) CreateMultipartUpload(ctx context.Context, key string) (string, error) {
	defer observe.Dependency("s3.create_multipart_upload")()

	resp, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}
	return *resp.UploadId, nil
}

func (s *S3Storage) GeneratePresignedPartURL(ctx context.Context, key string, uploadID string, partNumber int32) (string, error) {
	defer observe.Dependency("s3.presign_upload_part")()

	req, err := s.presigner.PresignUploadPart(ctx, &s3.UploadPartInput{
		Bucket:     &s.bucketName,
		Key:        &key,
		UploadId:   &uploadID,
		PartNumber: &partNumber,
	}, s3.WithPresignExpires(time.Duration(config.App.Upload.PresignExpirty)*time.Minute))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *S3Storage) CompleteMultipartUpload(ctx context.Context, key string, uploadID string, parts []domain.Chunk) error {
	defer observe.Dependency("s3.complete_multipart_upload")()

	var completedParts []types.CompletedPart

	for _, p := range parts {
		etag := p.ETag
		chunkNumber := int32(p.ChunkNumber)
		completedParts = append(completedParts, types.CompletedPart{
			ETag:       &etag,
			PartNumber: &chunkNumber,
		})
	}

	_, err := s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   &s.bucketName,
		Key:      &key,
		UploadId: &uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	return err
}

func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	defer observe.Dependency("s3.get_object")()
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
