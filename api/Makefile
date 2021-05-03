api_S3_bucket_name = bscfeesapi
bsc_fees_api_zip = bsc-fees.zip

help:
	@echo "--------- Available make commands ---------"
	@echo "-------------------------------------------"
	@echo "build                Builds the application"
	@echo "zip                  Zips the application. Useful when you want to push it to S3"
	@echo "sam                  Spins up a SAM instance with your app in it. Useful for local development"
	@echo "s3-api               Sends your app to the S3 bucket. Make sure that you have aws-cli installed"
	@echo "lambda               Updates Lambda Function code. Make sure that you have aws-cli installed"

.PHONY: build

build: 
	GOOS=linux GOARCH=amd64 go build -o bsc-fees ./cmd/bscfees

zip: build
	zip $(bsc_fees_api_zip) bsc-fees config.toml

s3-api: zip
	aws s3 cp bsc-fees.zip s3://$(api_S3_bucket_name)/$(bsc_fees_api_zip)

lambda: s3-api
	aws lambda update-function-code --function-name bsc-fees --s3-bucket $(api_S3_bucket_name) --s3-key $(bsc_fees_api_zip)

sam: build
	sam local start-api

clean:
	rm bsc-fees
