package vpc

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)


var cfgTest = setCfgTest()

func init() {
	cfg = cfgTest
	err := createDefaultVpc()
	if err != nil {
		log.Error(err)
	}
}

func setCfgTest() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error(err)
	}
	return cfg
}

func createDefaultVpc() error {
	svc := ec2.NewFromConfig(cfg)

	_, err := svc.CreateDefaultVpc(context.Background(), &ec2.CreateDefaultVpcInput{})
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func TestVpcDefaultId(t *testing.T) {

	got := VpcDefaultId(cfg)

	if assert.Contains(t, got, "vpc-") {
		log.Println("Default VPC ID:", got)
	}
}

func TestSubnetsIds(t *testing.T) {
	vpcDefaultId := VpcDefaultId(cfg)
	ids := SubnetsIds(vpcDefaultId, cfg)
	if assert.Equal(t, vpcDefaultId, *ids[0].VpcId) {
		for k, v := range ids {
			log.Println(k, *v.SubnetId)
		}
	}
}

func TestDeleteSubnets(t *testing.T) {
	vpcDefaultId := VpcDefaultId(cfg)
	ids := SubnetsIds(vpcDefaultId, cfg)

	err := DeleteSubnets(ids, cfg)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}

func TestIGWId(t *testing.T) {
	vpcDefaultId := VpcDefaultId(cfg)
	id := IGWId(vpcDefaultId, cfg)
	if assert.Contains(t, id, "igw-") {
		log.Println("IGW ID:", id)
	}
}

func TestDeleteIGW(t *testing.T) {
	vpcDefaultId := VpcDefaultId(cfg)
	id := IGWId(vpcDefaultId, cfg)

	err := DeleteIGW(id, vpcDefaultId, cfg)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}

func TestDeleteVpc(t *testing.T) {

	id := VpcDefaultId(cfg)

	err := DeleteVpc(id, cfg)
	if err != nil {
		log.Error(err)
		t.Fail()
	}

}

func TestDestroyVPC(t *testing.T) {
	err := createDefaultVpc()
	if err != nil {
		log.Error(err)
		t.Fail()
	}
	err = DestroyDefaultVpc(cfg)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}
