package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"time"
)

type Ec2Info struct {
	Name       *string
	InstanceId *string
	Status     *string
	Ip         *string
	Key        *string
}

func (c *Aws) CreateEc2(Ami string, Ec2Type string, Name string) (*Ec2Info, error) {
	svc := ec2.New(c.Sess)
	dateName := Name + time.Unix(time.Now().Unix(), 0).Format("_2006-01-02_15:04:05")
	keyRt, keyErr := svc.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: &dateName})
	if keyErr != nil {
		return nil, keyErr
	} //create ssh key
	runRt, runErr := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(Ami),
		InstanceType: aws.String(Ec2Type),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		KeyName:      &dateName,
	})
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
	})
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

func (c *Aws) ChangeEc2Ip(InstanceId string) (*string, error) {
	svc := ec2.New(c.Sess)
	desRt, desErr := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(InstanceId)},
			},
		},
	})
	if desErr != nil {
		return nil, desErr
	}
	if len(desRt.Addresses) != 0 {
		_, relErr := svc.ReleaseAddress(&ec2.ReleaseAddressInput{AllocationId: desRt.Addresses[0].AllocationId})
		if relErr != nil {
			return nil, relErr
		}
	}
	allRt, allErr := svc.AllocateAddress(&ec2.AllocateAddressInput{})
	if allErr != nil {
		return nil, allErr
	}
	_, assErr := svc.AssociateAddress(&ec2.AssociateAddressInput{
		AllocationId: allRt.AllocationId,
		InstanceId:   aws.String(InstanceId),
	})
	if assErr != nil {
		return nil, assErr
	}
	return allRt.PublicIp, nil
}

func (c *Aws) GetEc2Info(InstanceId string) (*Ec2Info, error) {
	svc := ec2.New(c.Sess)
	rt, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{aws.String(InstanceId)}})
	if err != nil {
		return nil, err
	}
	return &Ec2Info{
		Name:       rt.Reservations[0].Instances[0].Tags[0].Value,
		InstanceId: rt.Reservations[0].Instances[0].InstanceId,
		Status:     rt.Reservations[0].Instances[0].State.Name,
		Ip:         rt.Reservations[0].Instances[0].PublicIpAddress,
	}, nil
}

func (c *Aws) ListEc2() ([]*ec2.Reservation, error) {
	svc := ec2.New(c.Sess)
	rt, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{MaxResults: aws.Int64(100)})
	if err != nil {
		return nil, err
	}
	return rt.Reservations, nil
}

func (c *Aws) StartEc2(InstanceId string) error {
	svc := ec2.New(c.Sess)
	_, err := svc.StartInstances(&ec2.StartInstancesInput{InstanceIds: []*string{aws.String(InstanceId)}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Aws) StopEc2(InstanceId string) error {
	svc := ec2.New(c.Sess)
	_, err := svc.StopInstances(&ec2.StopInstancesInput{InstanceIds: []*string{aws.String(InstanceId)}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Aws) RebootEc2(InstanceId string) error {
	svc := ec2.New(c.Sess)
	_, err := svc.RebootInstances(&ec2.RebootInstancesInput{InstanceIds: []*string{aws.String(InstanceId)}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Aws) DeleteEc2(InstanceId string) error {
	svc := ec2.New(c.Sess)
	_, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{InstanceIds: []*string{aws.String(InstanceId)}})
	if err != nil {
		return err
	}
	return nil
}
