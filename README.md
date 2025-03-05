# cribl-storage-tool

A command-line tool for managing data across diverse storage solutions. As data lakes and data tiering become the norm, accessing data from multiple clouds and storage platforms can be challenging. The cribl-storage-tool helps you streamline these operations with ease.

## Table of Contents

- [Introduction](#introduction)
- [Why & How](#why--how)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
   - [S3 List Command](#s3-list-command)
   - [IAM Setup Command](#iam-setup-command)
- [Examples](#examples)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Introduction

The `cribl-storage-tool` is designed to simplify the management of data in heterogeneous storage environments. Whether you're working with S3, Azure Blob Storage, or other storage services, this tool provides a consistent interface to manage, list, and organize your data.

## Why & How

Managing data across multiple storage systems is increasingly complex. The cribl-storage-tool was built to:

- **Simplify Access:** Provide a unified command line interface to work with various storage solutions.
- **Improve Efficiency:** Automate common tasks such as listing, transferring, and managing storage objects.
- **Ensure Compatibility:** Support multiple cloud providers and on-premises storage solutions through extensible design.

## Features

- **Multi-Cloud Support:** Seamlessly integrate with AWS S3 and other storage platforms.
- **Command Line Interface:** Fast and efficient tool designed for automation and scripting.
- **Extensible Architecture:** Easily add support for new storage&data solutions as your needs evolve.

## Installation

Clone the repository and install the necessary dependencies:
Run this command in your terminal and it will install the binary file
```aiignore
./scripts/install.sh
```

## Usage

The tool is divided into multiple commands for interacting with different storage components. Below are examples for the two primary commands: listing S3 buckets and setting up IAM roles.

 - S3 List Command
```./cribl-storage-tool s3 list -h```
   - ```aiignore
      Usage:
      cribl-storage-tool s3 list [flags]
    
      Flags:
      -b, --bucket-file string   Path to JSON file containing S3 bucket names (optional)
      -f, --filter string        Filter bucket names containing the specified substring (optional)
      -h, --help                 help for list
      -o, --output string        Output format: text, json, or names (default "text")
      -p, --profile string       AWS profile to use for authentication (optional)
      -x, --regex string         Filter bucket names matching the specified regular expression (optional)
      -r, --region string        AWS region to target (optional)
   ```
 - IAM Setup Command Search it
   ```/cribl-storage-tool iam setup -h```
 - ```Usage:
   cribl-storage-tool iam setup [flags]
   
   Flags:
   -a, --account string            AWS Account ID to trust (required if --cribl-worker-arn not provided)
   -s, --action string             Action type for the IAM role (default: search) (default "search")
   -b, --bucket strings            Name of the S3 bucket to grant access (can specify multiple)
   -f, --bucket-file string        Path to JSON file containing S3 bucket names (optional)
   --cribl-worker-arn string   Cribl worker ARN (e.g., arn:aws:iam::ACCOUNT:role/WORKSPACE-WORKERGROUP)
   -e, --external-id string        External ID for the trust relationship (optional)
   -h, --help                      help for setup
   -p, --profile string            AWS profile to use for authentication (optional)
   -z, --region string             AWS region to target (optional)
   -r, --role string               Name of the IAM role to create or update (default "CrossAccountAccessRole")
   -g, --workergroup string        Worker group name (default: default) (default "default")
   -w, --workspace string          Workspace name (default: main) (default "main")
   ```
## Examples:
Lets go ahead and use my power account goatshipansible to list all the s3 buckets
```./cribl-storage-tool s3 list --profile goatshipansible```
```
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
Now that i have all my buckets listed i want to be able to search on Cribl Search for the buckets badcoffee and ckoamplifybucket
for account we will be using the cribl account trust relationship if you have that handy. In this case for my tenant: `47111295931415`
the workspace i'm going to link the bucket to is -w `contractors` .

Since I like pi im going to specify my external-id with the flag -e `31415` + the role name -r `elbcoffeee` here is the completed command for badcoffee bucket:

```./cribl-storage-tool iam setup --account 4711129531415 --profile goatshipansible --bucket badcoffee -e 314515 -r elbcoffee --workspace contractors ```

you will see a stream of logs and if successful `{"level":"info","command":"iam_setup","time":"2025-02-27T13:48:26-05:00","message":"IAM trust relationship setup completed successfully"}`

Speaking of stream lets go ahead and edit the command for stream to send data to the bucket from Stream:
```./cribl-storage-tool iam setup --account 4711129531415 --profile goatshipansible --bucket badcoffee -e 314515 -r elbcoffee --workspace contractors --workergroup default --action send```

- Using the filter || regex
```
```
./cribl-storage-tool s3 list --profile goatshipansible --filter lake
```aiignore
zamorofthat@29JH7X-pi cribl-storage-tool % ./cribl-storage-tool s3 list --profile goatshipansible --filter lake

Listing S3 Buckets:
 - aws-security-data-lake-us-east-1-55555
 - aws-security-data-lake-us-east-2-55555
 - aws-security-data-lake-us-west-1-55555
 - aws-security-data-lake-us-west-2-55555
 - seclake-customsource

```
./cribl-storage-tool s3 list --profile goatshipansible --regex "lake.*"\
```aiignore
zamorofthat@29JH7X-pi cribl-storage-tool % ./cribl-storage-tool s3 list --profile goatshipansible --regex "lake.*"

Listing S3 Buckets:
 - aws-security-data-lake-us-east-1-55555
 - aws-security-data-lake-us-east-2-55555
 - aws-security-data-lake-us-west-1-55555
 - aws-security-data-lake-us-west-2-55555
 - seclake-customsource
```
Create a bucket file to loop over in bash 
```aiignore
cribl-storage-tool s3 list -p goatshipansible -o names > goats.txt 
badcoffee
criblcoffeeroute53
criblcompetitorsbucket
kerno-samples-cdd9cc33-9d8c-442f-af10-6e07196e0d71
seclake-customsource

```