# Document service

The document service currently owns document uploads. Its primary flow uses S3 multipart uploads with presigned URLs: the service creates the upload, the client sends each chunk directly to S3, and the service completes the upload using the returned ETags.

## Run locally

1. Copy `.env.example` to `.env` and add the S3 bucket and credentials.
2. Start the service:

```bash
go run ./cmd/server
```

The default HTTP port is `8090`.

## Upload a document in chunks

Initialize the upload:

```bash
curl -X POST http://localhost:8090/documents/uploads/init \
  -H 'Content-Type: application/json' \
  -d '{"file_name":"document.pdf","file_size":12345678,"content_type":"application/pdf"}'
```

The response contains `upload_id`, `storage_key`, `chunk_size`, and one presigned URL per part. Upload each part directly to its URL using the matching `part_number`, then collect each S3 response `ETag` and complete the upload:

```bash
curl -X POST http://localhost:8090/documents/uploads/complete \
  -H 'Content-Type: application/json' \
  -d '{"upload_id":"...","storage_key":"...","parts":[{"part_number":1,"etag":"\"...\""}]}'
```

The client must upload parts in order, using the configured chunk size for every part except the final part.

The shorter `/documents/upload` and `/documents/complete-upload` routes are compatibility aliases for the same initialize and complete operations:

```bash
curl -X POST http://localhost:8090/documents/upload \
  -H 'Content-Type: application/json' \
  -d '{"file_name":"document.pdf","file_size":12345678,"content_type":"application/pdf"}'
```

Configuration:

- `AWS_BUCKET`: required bucket name, matching the existing `elevate` service convention. `S3_BUCKET` is also accepted for compatibility.
- `S3_PREFIX`: optional object key prefix, default `documents`.
- `S3_PRESIGN_EXPIRY_MINUTES`: presigned part URL lifetime, default `15` minutes.
- `S3_CHUNK_SIZE_MB`: part size, default `8` MB.
- `S3_MAX_FILE_SIZE_MB`: request limit, default `512` MB.
- `S3_ENDPOINT` and `S3_FORCE_PATH_STYLE`: optional S3-compatible storage settings.
