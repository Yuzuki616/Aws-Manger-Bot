package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"time"
)

func (p *Aws) CreateWl(Zone string, VpcId string) error {
	svc := ec2.New(p.Sess)
	_, subErr := svc.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZone: aws.String(Zone),
		CidrBlock:        aws.String("172.31.256.0/20"),
		VpcId:            aws.String(VpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("subnet"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("name"),
						Value: aws.String("aws_manger_subnet"),
					},
				},
			},
		},
	})
	if subErr != nil {
		return subErr
	}
	ca, caErr := svc.CreateCarrierGateway(&ec2.CreateCarrierGatewayInput{
		VpcId: aws.String(VpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("carrier-gateway"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("name"),
						Value: aws.String("aws_manger_gateway"),
					},
				},
			},
		},
	})
	if caErr != nil {
		return caErr
	}
	route, routeErr := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(VpcId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("route-table"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("name"),
						Value: aws.String("aws_manger_route"),
					},
				},
			},
		},
	})
	if routeErr != nil {
		return routeErr
	}
	_, assErr := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: route.RouteTable.RouteTableId,
		GatewayId:    ca.CarrierGateway.CarrierGatewayId,
	})
	if assErr != nil {
		return assErr
	}
	return nil
}

func (p Aws) GetGatewayInfo() (*ec2.DescribeCarrierGatewaysOutput, error) {
	svc := ec2.New(p.Sess)
	ca, err := svc.DescribeCarrierGateways(&ec2.DescribeCarrierGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String("aws_manger_gateway")},
			},
		}})
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (p *Aws) CreateEc2Wl(SubId string, Ami string, Name string) (*Ec2Info, error) {
	svc := ec2.New(p.Sess)
	dateName := Name + time.Unix(time.Now().Unix(), 0).Format("_2006-01-02_15:04:05")
	keyRt, keyErr := svc.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: &dateName})
	if keyErr != nil {
		return nil, keyErr
	} //创建ssh密钥
	secRt, secErr := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(dateName + "security"),
		Description: aws.String("A security group for aws manger bot"),
	}) //创建安全组
	if secErr != nil {
		return nil, secErr
	}
	_, authSecInErr := svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: secRt.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
				FromPort: aws.Int64(-1),
				ToPort:   aws.Int64(-1),
			},
		},
	}) //添加入站规则
	if authSecInErr != nil {
		return nil, authSecInErr
	}
	runRt, runErr := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(Ami),
		InstanceType:     aws.String("t2.medium"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          &dateName,
		SecurityGroupIds: []*string{secRt.GroupId},
		SubnetId:         aws.String(SubId),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{{
			AssociateCarrierIpAddress: aws.Bool(true),
		}},
	}) //创建ec2实例
	if runErr != nil {
		return nil, runErr
	}
	_, tagErr := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runRt.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(Name),
			},
		},
	}) //创建标签
	if tagErr != nil {
		return nil, tagErr
	}
	return &Ec2Info{
		Name:       &Name,
		InstanceId: runRt.Instances[0].InstanceId,
		Status:     runRt.Instances[0].State.Name,
		Key:        keyRt.KeyMaterial,
	}, nil
}
