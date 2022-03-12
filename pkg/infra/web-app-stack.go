package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type WebAppStackProps struct {
	awscdk.StackProps
	EnvName string
}

func NewWebAppStack(scope constructs.Construct, id string, props *WebAppStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	stackName := props.StackName

	// S3 Bucket
	bucketPartialName := "bucket" // Bucket name will be generated from ID which already includes stack details
	bucket := awss3.NewBucket(stack, &bucketPartialName, &awss3.BucketProps{
		WebsiteIndexDocument: jsii.String("index.html"),
		PublicReadAccess:     jsii.Bool(true),
		Encryption:           awss3.BucketEncryption_S3_MANAGED,
		EnforceSSL:           jsii.Bool(true),
		Versioned:            jsii.Bool(false),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedOrigins: jsii.Strings("*"),
				AllowedMethods: &[]awss3.HttpMethods{awss3.HttpMethods_GET},
				AllowedHeaders: jsii.Strings("*"),
			},
		},
	})
	NewInfraParameter(stack, props.EnvName, ParamWebAppBucketName, *bucket.BucketName())

	// S3 Origin Access Identity
	oiaName := *stackName + "-origin-access-identity"
	oia := awscloudfront.NewOriginAccessIdentity(stack, &oiaName, nil)
	bucket.GrantRead(oia, nil)

	// CloudFront
	cdnName := *stackName + "-cloudfront-distribution"
	cdn := awscloudfront.NewCloudFrontWebDistribution(stack, &cdnName, &awscloudfront.CloudFrontWebDistributionProps{
		OriginConfigs: &[]*awscloudfront.SourceConfiguration{
			{
				S3OriginSource: &awscloudfront.S3OriginConfig{
					S3BucketSource:       bucket,
					OriginAccessIdentity: oia,
				},
				Behaviors: &[]*awscloudfront.Behavior{
					{
						IsDefaultBehavior: jsii.Bool(true),
					},
				},
			},
		},
	})
	NewInfraParameter(stack, props.EnvName, ParamWebAppDistributionId, *cdn.DistributionId())

	return stack
}
