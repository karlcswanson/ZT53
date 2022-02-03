package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	_ "github.com/joho/godotenv/autoload"
	ztcentral "github.com/zerotier/go-ztcentral"
)

type Member struct {
	name string
	ip   string
}

type Network_Members struct {
	members []Member
}

func main() {
	session := r53Session()
	members := getNetworkDevices(os.Getenv("ZT_NETWORK"))
	for _, m := range members.members {
		fmt.Printf("%s\n", m.name)
	}
	changes := changeList(members)

	updateR53(session, changes)
}

func r53Session() *route53.Route53 {
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)
	return svc
}

func changeList(member_list Network_Members) route53.ChangeBatch {
	var changes route53.ChangeBatch
	for _, m := range member_list.members {
		hostname := m.name + "." + os.Getenv("DOMAIN")
		c := route53.Change{
			Action: aws.String("UPSERT"),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: aws.String(hostname),
				Type: aws.String("A"),
				ResourceRecords: []*route53.ResourceRecord{
					{
						Value: aws.String(string(m.ip)),
					},
				},
				TTL: aws.Int64(60),
			},
		}
		changes.Changes = append(changes.Changes, &c)
	}
	return changes
}

func updateR53(svc *route53.Route53, changeList route53.ChangeBatch) {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  &changeList,
		HostedZoneId: aws.String(os.Getenv("R53_ZONE")),
	}
	fmt.Printf("%s\n", params)
	resp, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Change Response:")
	fmt.Println(resp)

}

func getNetworkDevices(network_id string) Network_Members {
	c, err := ztcentral.NewClient(os.Getenv("ZT_TOKEN"))

	var member_list Network_Members

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	members, err := c.GetMembers(ctx, network_id)
	if err != nil {
		log.Println("error:", err.Error())
		os.Exit(1)
	}

	for _, m := range members {
		log.Printf("\t%s\t %s %s", *m.Id, *m.Name, *m.Config.IpAssignments)
		ipa := *m.Config.IpAssignments
		member := Member{name: *m.Name, ip: ipa[0]}
		member_list.members = append(member_list.members, member)
	}

	return member_list
}
