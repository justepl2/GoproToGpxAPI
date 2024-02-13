run:
	docker-compose up -d
	cd src && go mod tidy
	cd src && go run main.go
	# aws --endpoint-url=http://localhost:4566 s3api create-bucket --bucket gopro-to-gpx-api --create-bucket-configuration LocationConstraint=eu-west-3
	# aws --endpoint-url=http://localhost:4566 s3api list-buckets
