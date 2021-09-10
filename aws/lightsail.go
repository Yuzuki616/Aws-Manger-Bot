package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"time"
)

type LsInfo struct {
	Name   *string
	Ip     *string
	IpName *string
	Status *string
	Key    *string
}

func (c *Aws) CreateLs(Name string, Region string, BlueId string, BundleId string) (*LsInfo, error) {
	svc := lightsail.New(c.Sess)
	dateName := Name + time.Unix(time.Now().Unix(), 0).Format("_2006-01-02_15:04:05")
	key, keyErr := svc.CreateKeyPair(&lightsail.CreateKeyPairInput{KeyPairName: aws.String(dateName)})
	if keyErr != nil {
		return nil, keyErr
	}
	lsRt, lsErr := svc.CreateInstances(&lightsail.CreateInstancesInput{
		AvailabilityZone: aws.String(Region),
		BlueprintId:      aws.String(BlueId),
		BundleId:         aws.String(BundleId),
		InstanceNames:    []*string{aws.String(Name)},
		KeyPairName:      aws.String(dateName),
		Tags: []*lightsail.Tag{
			{
				Key:   aws.String("SIp"),
				Value: aws.String(dateName),
			},
		},
	})
	if lsErr != nil {
		return nil, lsErr
	}
	_, allErr := svc.AllocateStaticIp(&lightsail.AllocateStaticIpInput{StaticIpName: aws.String(dateName)})
	if allErr != nil {
		return nil, allErr
	}
	_, attErr := svc.AttachStaticIp(&lightsail.AttachStaticIpInput{
		StaticIpName: aws.String(dateName),
		InstanceName: aws.String(Name),
	})
	if attErr != nil {
		return nil, attErr
	}
	return &LsInfo{
		Name:   &Name,
		Status: lsRt.Operations[0].Status,
		Key:    key.PrivateKeyBase64,
	}, nil
}

func (c *Aws) GetLsInfo(Name string) (*LsInfo, error) {
	svc := lightsail.New(c.Sess)
	rt, err := svc.GetInstance(&lightsail.GetInstanceInput{InstanceName: aws.String(Name)})
	if err != nil {
		return nil, err
	}
	return &LsInfo{
		Name:   rt.Instance.Name,
		Ip:     rt.Instance.PublicIpAddress,
		IpName: rt.Instance.Tags[0].Value,
		Status: rt.Instance.State.Name,
	}, nil
}

func (c *Aws) ListLs() ([]*lightsail.Instance, error) {
	svc := lightsail.New(c.Sess)
	rt, err := svc.GetInstances(&lightsail.GetInstancesInput{})
	if err != nil {
		return nil, err
	}
	return rt.Instances, nil
}

func (c *Aws) ChangeLsIp(Name string) error {
	getRt, getErr := c.GetLsInfo(Name)
	if getErr != nil {
		return getErr
	}
	svc := lightsail.New(c.Sess)
	_, detErr := svc.DetachStaticIp(&lightsail.DetachStaticIpInput{StaticIpName: getRt.IpName})
	if detErr != nil {
		return detErr
	}
	_, relErr := svc.ReleaseStaticIp(&lightsail.ReleaseStaticIpInput{StaticIpName: getRt.IpName})
	if relErr != nil {
		return relErr
	}
	_, allErr := svc.AllocateStaticIp(&lightsail.AllocateStaticIpInput{StaticIpName: getRt.IpName})
	if allErr != nil {
		return allErr
	}
	_, attErr := svc.AttachStaticIp(&lightsail.AttachStaticIpInput{
		StaticIpName: getRt.IpName,
		InstanceName: aws.String(Name),
	})
	if attErr != nil {
		return attErr
	}
	return nil
}

func (c *Aws) StopLs(Name string) error {
	svc := lightsail.New(c.Sess)
	_, err := svc.StopInstance(&lightsail.StopInstanceInput{InstanceName: aws.String(Name)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Aws) StartLs(Name string) error {
	svc := lightsail.New(c.Sess)
	_, err := svc.StartInstance(&lightsail.StartInstanceInput{InstanceName: aws.String(Name)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Aws) DeleteLs(Name string) error {
	svc := lightsail.New(c.Sess)
	_, err := svc.DeleteInstance(&lightsail.DeleteInstanceInput{InstanceName: aws.String(Name)})
	if err != nil {
		return err
	}
	return nil
}
