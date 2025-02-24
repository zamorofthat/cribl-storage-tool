# cribl-storage-tool
./cribl-storage-tool s3 list --profile goatshipansible
```aiignore
azamora@29JH7X-luQT cribl-storage-tool % ./cribl-storage-tool s3 list --profile goatshipansible

Listing S3 Buckets:
 - aws-cloudtrail-logs-55555-55555
 - aws-security-data-lake-us-east-1-55555
 - aws-security-data-lake-us-east-2-55555
 - aws-security-data-lake-us-west-1-55555
 - aws-security-data-lake-us-west-2-55555
 - badcoffee
 - ckoamplifybucket
 - ckoamplifybucketparquet
 - config-bucket-55555
 - criblcoffeeroute53
 - 55555-55555-55555-9d8c-442f-af10-55555
 - o11ys3bucket
 - seclake-customsource

```
./cribl-storage-tool s3 list --profile goatshipansible --filter lake
```aiignore
azamora@29JH7X-luQT cribl-storage-tool % ./cribl-storage-tool s3 list --profile goatshipansible --filter lake

Listing S3 Buckets:
 - aws-security-data-lake-us-east-1-55555
 - aws-security-data-lake-us-east-2-55555
 - aws-security-data-lake-us-west-1-55555
 - aws-security-data-lake-us-west-2-55555
 - seclake-customsource

```

./cribl-storage-tool s3 list --profile goatshipansible --regex "lake.*"\
```aiignore
azamora@29JH7X-luQT cribl-storage-tool % ./cribl-storage-tool s3 list --profile goatshipansible --regex "lake.*"

Listing S3 Buckets:
 - aws-security-data-lake-us-east-1-55555
 - aws-security-data-lake-us-east-2-55555
 - aws-security-data-lake-us-west-1-55555
 - aws-security-data-lake-us-west-2-55555
 - seclake-customsource
477358655677 
```
./cribl-storage-tool iam setup --account 47731415 --profile criblcoffee --region us-east-1 --bucket badcoffee -e 314515 -r elbcoffee --workspace keynote
```aiignore
{"level":"info","command":"iam_setup","component":"iam_client","role_name":"elbcoffee","trusted_account_id":"477358655677","workspace":"keynote","workergroup":"default","action":"search","bucket_names":["badcoffee"],"time":"2025-02-19T18:21:56-05:00","message":"setting up trust relationship"}
{"level":"debug","command":"iam_setup","component":"iam_client","role_name":"elbcoffee","trusted_account_id":"477358655677","workspace":"keynote","workergroup":"default","bucket_names":["badcoffee"],"time":"2025-02-19T18:21:56-05:00","message":"validating inputs"}
{"level":"debug","command":"iam_setup","component":"iam_client","role_name":"elbcoffee","trusted_account_id":"477358655677","workspace":"keynote","workergroup":"default","bucket_names":["badcoffee"],"time":"2025-02-19T18:21:56-05:00","message":"input validation successful"}

```