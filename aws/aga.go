package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	aga "github.com/aws/aws-sdk-go/service/globalaccelerator"
	"time"
)

type AgaInfo struct {
	Name     *string
	Status   *string
	Ip       []*aga.IpSet
	Protocol *string
	Port     []*aga.PortOverride
}

func (c *Aws) CreateAga(Name string, Region string, InstanceId string, ListenerArn string, TargetPort int64, SourcePort int64) (*AgaInfo, error) {
	svc := aga.New(c.Sess)
	IdempotencyToken := time.Unix(time.Now().Unix(), 0).Format("2006-01-02_15:04:05")
	createAccRt, createAccErr := svc.CreateAccelerator(&aga.CreateAcceleratorInput{
		Name:             aws.String(Name),
		Enabled:          aws.Bool(false),
		IdempotencyToken: aws.String(IdempotencyToken),
	})
	if createAccErr != nil {
		return nil, createAccErr
	}
	createLiRt, createLiErr := svc.CreateListener(&aga.CreateListenerInput{
		AcceleratorArn: createAccRt.Accelerator.AcceleratorArn,
		PortRanges: []*aga.PortRange{
			{
				FromPort: aws.Int64(1),
				ToPort:   aws.Int64(65535),
			},
		},
		Protocol: aws.String("TCP"),
	})
	if createLiErr != nil {
		return nil, createLiErr
	}
	createEndRt, createEndErr := svc.CreateEndpointGroup(&aga.CreateEndpointGroupInput{
		EndpointGroupRegion: aws.String(Region),
		IdempotencyToken:    aws.String(IdempotencyToken),
		ListenerArn:         aws.String(ListenerArn),
		HealthCheckPort:     aws.Int64(22),
		EndpointConfigurations: []*aga.EndpointConfiguration{
			{
				EndpointId: aws.String(InstanceId),
			},
		},
		PortOverrides: []*aga.PortOverride{
			{
				EndpointPort: aws.Int64(TargetPort),
				ListenerPort: aws.Int64(SourcePort),
			},
		},
	})
	if createEndErr != nil {
		return nil, createEndErr
	}
	return &AgaInfo{
		Name:     createAccRt.Accelerator.Name,
		Status:   createAccRt.Accelerator.Status,
		Ip:       createAccRt.Accelerator.IpSets,
		Protocol: createLiRt.Listener.Protocol,
		Port:     createEndRt.EndpointGroup.PortOverrides,
	}, nil
}

func (c *Aws) ListAga() ([]*aga.Accelerator, error) {
	svc := aga.New(c.Sess)
	rt, err := svc.ListAccelerators(&aga.ListAcceleratorsInput{})
	if err != nil {
		return nil, err
	}
	return rt.Accelerators, err
}

func (c *Aws) GetAgaInfo(AcceleratorArn string) (*AgaInfo, error) {
	svc := aga.New(c.Sess)
	accRt, accErr := svc.DescribeAccelerator(&aga.DescribeAcceleratorInput{AcceleratorArn: aws.String(AcceleratorArn)})
	if accErr != nil {
		return nil, accErr
	}
	liRt, liErr := svc.ListListeners(&aga.ListListenersInput{AcceleratorArn: accRt.Accelerator.AcceleratorArn})
	if liErr != nil {
		return nil, liErr
	}
	endRt, endErr := svc.ListEndpointGroups(&aga.ListEndpointGroupsInput{ListenerArn: liRt.Listeners[0].ListenerArn})
	if endErr != nil {
		return nil, endErr
	}
	return &AgaInfo{
		Name:     accRt.Accelerator.Name,
		Status:   accRt.Accelerator.Status,
		Ip:       accRt.Accelerator.IpSets,
		Protocol: liRt.Listeners[0].Protocol,
		Port:     endRt.EndpointGroups[0].PortOverrides,
	}, nil
}
