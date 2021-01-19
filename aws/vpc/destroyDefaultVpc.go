package vpc

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	log "github.com/sirupsen/logrus"
)

var cfg = setCfg()
var region = "us-east-1"

func setCfg() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Error(err)
	}
	return cfg
}

//

func VpcDefaultId(cfg aws.Config) string {

	svc := ec2.NewFromConfig(cfg)
	req, err := svc.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("isDefault"),
				Values: []string{"true"},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	return *req.Vpcs[0].VpcId // Since the account is New, there only should be ONE DEFAULT VPC
}

func SubnetsIds(id string, cfg aws.Config) []types.Subnet {
	svc := ec2.NewFromConfig(cfg)
	req, err := svc.DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{id},
			},
		},
	})

	if err != nil {
		log.Error(err)
		return nil
	}

	return req.Subnets
}

func DeleteSubnets(ids []types.Subnet, cfg aws.Config) error {
	svc := ec2.NewFromConfig(cfg)

	for _, v := range ids {
		_, err := svc.DeleteSubnet(context.Background(), &ec2.DeleteSubnetInput{
			SubnetId: aws.String(*v.SubnetId),
			DryRun:   false,
		})
		if err != nil {
			log.Error(err)
			return err
		} else {
			log.Println("Deleted: ", *v.SubnetId)
		}

	}
	return nil
}

func IGWId(id string, cfg aws.Config) string {
	svc := ec2.NewFromConfig(cfg)

	req, err := svc.DescribeInternetGateways(context.Background(), &ec2.DescribeInternetGatewaysInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []string{id},
			},
		},
	})

	if err != nil {
		log.Error(err)
		return ""
	}

	return *req.InternetGateways[0].InternetGatewayId

}

func DeleteIGW(igwId string, vpcId string, cfg aws.Config) error {

	svc := ec2.NewFromConfig(cfg)

	_, err :=  svc.DetachInternetGateway(context.Background(), &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	})
	if err != nil{
		log.Error(err)
		return err
	}
	log.Println("Dettached IGW:", igwId, "from VPC:", vpcId)
	_ , err = svc.DeleteInternetGateway(context.Background(), &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
	})
	if err != nil{
		log.Error(err)
		return err
	}
	log.Println("Deleted: ", igwId)
	return nil
}

func DeleteVpc(id string, cfg aws.Config) error {

	svc := ec2.NewFromConfig(cfg)

	_, err := svc.DeleteVpc(context.Background(), &ec2.DeleteVpcInput{
		VpcId: aws.String(id),
	})

	if err != nil {
		log.Error(err)
		return err
	}
	log.Println("Deleted:", id)
	return nil
}

func DestroyDefaultVpc(cfg aws.Config) error {
	// Get Default VPC ID
	vpcID := VpcDefaultId(cfg)

	// Delete all default resources
	err := DeleteSubnets(SubnetsIds(vpcID,cfg),cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Println("Deleted all the default subnets")
	err = DeleteIGW(IGWId(vpcID,cfg),vpcID,cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Println("Deleted the default IGW")

	// Delete the Default VPC
	err = DeleteVpc(vpcID,cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Println("Deleted the default VPC")

	return nil
}