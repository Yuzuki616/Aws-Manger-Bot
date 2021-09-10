package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Aws struct {
	Sess *session.Session
}

func New(Region string, Id string, Secret string) (*Aws, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(Id, Secret, ""),
	})
	c := &Aws{
		Sess: sess,
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
