package main

import (
	"fmt"
	"time"
	"context"
	"flag"
	"log"
	"google.golang.org/grpc"
	v1 "secondwork/pubsub/v1"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/genproto/protobuf/field_mask"

)

var(
	addr = flag.String("addr", "127.0.0.1:1236", "the address to connect to")

	projectIdTest string
)

func main() {
	flag.Parse()
	projectIdTest="PCG-HOMEWORK2"
	conn,err :=grpc.Dial(*addr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err !=nil{
		log.Printf("fail to dial :%v",err)
		return 
	}
	fmt.Printf("Client booting successful!\n")

	defer conn.Close()
	client:=v1.NewPublisherClient(conn)

	fmt.Printf("start CreateTopic\n")
	ctx, e := context.WithTimeout(context.Background(), time.Second*10)
	defer e()
	topic1,err1:=client.CreateTopic(ctx,&v1.Topic{Name: getTotalNameByShortName("bbb", "topic")})
	if err1 !=nil{
		log.Fatalf("client create topic1 falied!")
	}else {
		fmt.Printf("create topic1 success : %s\n",topic1.Name)
	}

	ctx2, e2 := context.WithTimeout(context.Background(), time.Second*10)
	defer e2()
	topic2,err2:=client.CreateTopic(ctx2,&v1.Topic{Name: getTotalNameByShortName("aaa", "topic")})
	if err2 !=nil{
		log.Fatalf("client create topic2 falied!")
	}else {
		fmt.Printf("create topic2 success : %s\n",topic2.Name)
	}
	fmt.Printf("end CreateTopic\n")

	fmt.Printf("start UpdateTopic\n")
	ctx, e = context.WithTimeout(context.Background(), time.Second*10)
	update_r:=&v1.UpdateTopicRequest{
		Topic:topic1,
		UpdateMask:&field_mask.FieldMask{Paths :[]string{"KmsKetName"}},
	}
	defer e()
	_,errx:=client.UpdateTopic(ctx,update_r)
	if errx !=nil{
		log.Fatalf("error when update  %s\n",topic1.Name)
	}
	fmt.Printf("end UpdateTopic\n")

	fmt.Printf("start Publish\n")
	ctx,e=context.WithTimeout(context.Background(),time.Second*10)
	mess_l:=make([]*v1.PubsubMessage,2)
	r:=&v1.PublishRequest{
		Topic:topic1.Name,
		Messages:mess_l,
	}
	defer e()
	response,errr:=client.Publish(ctx,r)
	if errr!=nil{
		log.Fatalf("Publish error!\n")
	}
	for _,res:=range response.MessageIds{
		fmt.Printf("message id :%s\n",res)
	}
	fmt.Printf("end Publish\n")

	fmt.Printf("start GetTopic\n")
	ctx,e=context.WithTimeout(context.Background(),time.Second*10)
	Get_r:=&v1.GetTopicRequest{Topic:topic1.Name}
	defer e()
	Get_topic,err_get:=client.GetTopic(ctx,Get_r)
	if err_get!=nil{
		log.Fatalf("error GetTopic %s\n",Get_topic.Name)
	}	
	fmt.Printf("Get topic %s\n",Get_topic.Name)
	fmt.Printf("end GetTopic\n")

	fmt.Printf("start ListTopics\n")

	ctx,e=context.WithTimeout(context.Background(),time.Second*10)
	List_r:=&v1.ListTopicsRequest{}
	defer e()
	List_topic,err_lis:=client.ListTopics(ctx,List_r)
	if err_lis !=nil{
		log.Fatalf("error ListTopics\n")
	}
	for _,t_top:=range List_topic.Topics {
		fmt.Printf("List topic %s\n",t_top.Name)
	}
	fmt.Printf("end ListTopics\n")
	
	fmt.Printf("start DeleteTopic\n")
	ctx,e=context.WithTimeout(context.Background(),time.Second*10)
	Deltop_r:=&v1.DeleteTopicRequest{
		Topic:topic1.Name}
	defer e()
	_,errDel:=client.DeleteTopic(ctx,Deltop_r)
	if errDel !=nil{
		log.Fatalf("error DeleteTopic %s\n",topic1.Name)
	}
	ctx,e=context.WithTimeout(context.Background(),time.Second*10)
	List_r=&v1.ListTopicsRequest{}
	defer e()
	List_topic,err_lis=client.ListTopics(ctx,List_r)
	if err_lis !=nil{
		log.Fatalf("error ListTopics\n")
	}
	for _,t_top:=range List_topic.Topics {
		fmt.Printf("List topic %s\n",t_top.Name)
	}
	fmt.Printf("end DeleteTopic\n")
	

}
func getTotalNameByShortName(name string, type_ string) string {
	switch type_ {
	case "schema":
		return "projects/" + projectIdTest + "/schemas/" + name
	case "topic":
		return "projects/" + projectIdTest + "/topics/" + name
	case "subscription":
		return "projects/" + projectIdTest + "/subscriptions/" + name
	default:
		return ""
	}
}
