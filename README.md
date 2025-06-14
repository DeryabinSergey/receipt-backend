Simple deploy from source 
```bash
gcloud run deploy receipt-backend --region=europe-west3 --allow-unauthenticated --source .
```

```bash
GOOS=linux GOARCH=amd64 go build -o bin/app cmd/app/main.go 

docker buildx build --platform linux/amd64 -t gcr.io/money-advice-462707/receipt-backend -f bin/Dockerfile .
docker push gcr.io/money-advice-462707/receipt-backend
```