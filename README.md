# cribl-storage-tool

A command-line tool for managing data across diverse storage solutions. As data lakes and data tiering become the norm, accessing data from multiple clouds and storage platforms can be challenging. The cribl-storage-tool helps you streamline these operations with ease.

## Table of Contents

- [Introduction](#introduction)
- [Why & How](#why--how)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
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


## Usage

The tool is divided into multiple commands for interacting with different storage components. Below are examples for the two primary commands: listing S3 buckets and setting up IAM roles.

 - S3 List Command
```./cribl-storage-tool s3 list [flags]```
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
```
The s3 list command helps you retrieve a list of S3 buckets (or objects) based on various filters and output formats.
```bash
git clone https://github.com/yourusername/cribl-storage-tool.git
cd cribl-storage-tool
# Follow further installation instructions specific to your environment
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
