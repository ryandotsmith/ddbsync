package aws

type endpointStruct struct {
	Version   int
	Endpoints map[string]endpointEntry
}

type endpointEntry struct {
	Endpoint string
}

var endpointsMap = endpointStruct{
	Version: 2,
	Endpoints: map[string]endpointEntry{
		"*/*": endpointEntry{
			Endpoint: "{service}.{region}.amazonaws.com",
		},
		"*/cloudfront": endpointEntry{
			Endpoint: "cloudfront.amazonaws.com",
		},
		"*/iam": endpointEntry{
			Endpoint: "iam.amazonaws.com",
		},
		"*/importexport": endpointEntry{
			Endpoint: "importexport.amazonaws.com",
		},
		"*/route53": endpointEntry{
			Endpoint: "route53.amazonaws.com",
		},
		"*/sts": endpointEntry{
			Endpoint: "sts.amazonaws.com",
		},
		"ap-northeast-1/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"ap-southeast-1/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"ap-southeast-2/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"cn-north-1/*": endpointEntry{
			Endpoint: "{service}.{region}.amazonaws.com.cn",
		},
		"eu-central-1/s3": endpointEntry{
			Endpoint: "{service}.{region}.amazonaws.com",
		},
		"eu-west-1/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"sa-east-1/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"us-east-1/s3": endpointEntry{
			Endpoint: "s3.amazonaws.com",
		},
		"us-east-1/sdb": endpointEntry{
			Endpoint: "sdb.amazonaws.com",
		},
		"us-gov-west-1/iam": endpointEntry{
			Endpoint: "iam.us-gov.amazonaws.com",
		},
		"us-gov-west-1/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"us-gov-west-1/sts": endpointEntry{
			Endpoint: "sts.us-gov-west-1.amazonaws.com",
		},
		"us-west-1/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
		"us-west-2/s3": endpointEntry{
			Endpoint: "s3-{region}.amazonaws.com",
		},
	},
}
