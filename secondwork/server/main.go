package main

import (
	"fmt"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	v1 "secondwork/pubsub/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

)

var(
	topic_list map[string]*Topic
	message_list map[string]*v1.ReceivedMessage
)

type PubServer struct{
	v1.UnimplementedPublisherServer
	ack_id int
}

type Topic struct {
	data *v1.Topic
	subscri map[string]*Subscription
}

type Subscription struct {
	topic *Topic
	data  *v1.Subscription
	mess_list map[string]*v1.ReceivedMessage
}

func (s *PubServer) CreateTopic(ctx context.Context, t *v1.Topic)(*v1.Topic,error){
	topic,e:=topic_list[t.Name]
	if  e!=true{
		new_topic:=&v1.Topic{
		Name:t.Name,
	}
	topic_list[t.Name]=&Topic{
		data:new_topic,
		subscri: map[string]*Subscription{},
	}

	fmt.Println("topic：",t.Name,"successfully created!")
		return new_topic, nil
	}else{
		return topic.data,nil
	    }
}

func (s *PubServer) UpdateTopic(ctx context.Context, r *v1.UpdateTopicRequest)(*v1.Topic,error){
	topic,e:=topic_list[r.Topic.Name]
	if e !=true {
		log.Fatalf("Failed when UpdateTopic %s\n",r.Topic.Name)
		return nil,nil
	}
	for _,path:=range r.UpdateMask.Paths {
		switch path{
		case "Labels": topic.data.Labels=r.Topic.Labels
	case "MessageStorgePolicy": topic.data.MessageStoragePolicy=r.Topic.MessageStoragePolicy
	case"KmsKetName" :topic.data.KmsKeyName=r.Topic.KmsKeyName
	case"SchemaSettings":topic.data.SchemaSettings=r.Topic.SchemaSettings
	case"MessageRetentionDuration" :topic.data.MessageRetentionDuration=r.Topic.MessageRetentionDuration
}}
	return topic.data,nil
}

func (s *PubServer) Publish(ctx context.Context, p *v1.PublishRequest)(*v1.PublishResponse,error){
	topic,e:=topic_list[p.Topic]
	if e!=true{
		log.Fatalf("Publish failed when looking for list\n")
		return nil,nil
	}
	var mess_ids []string
	for _,mess:=range p.Messages{
		id:=s.ack_id
		s.ack_id++
		mess_ids=append(mess_ids,fmt.Sprintf("%d",id))
		tmp:=&v1.ReceivedMessage{
			AckId:fmt.Sprintf("%d",id),
			Message:mess,
		}
		for _,sub:=range topic.subscri{
			sub.mess_list[fmt.Sprintf("%d",id)]=tmp	
		}
	}
	return &v1.PublishResponse{MessageIds:mess_ids} ,nil

}

func (s *PubServer) GetTopic(ctx context.Context , r *v1.GetTopicRequest)(*v1.Topic, error){
	topic,err:=topic_list[r.Topic]
	if err!=true{
		log.Fatalf("Can't get Topic %s\n",r.Topic)
		return nil,nil
	}
	return topic.data,nil
}

func (s *PubServer) ListTopics(ctx context.Context, r *v1.ListTopicsRequest)(*v1.ListTopicsResponse,error){
	topics:=make([]*v1.Topic,len(topic_list))
	index:=0
	for _,t:=range topic_list {
		topics[index]=t.data
		index++
	}
	return &v1.ListTopicsResponse{Topics:topics},nil
}

func (s *PubServer) ListTopicSubscriptions(ctx context.Context, r *v1.ListTopicSubscriptionsRequest)(*v1.ListTopicSubscriptionsResponse,error){
	topic,e:=topic_list[r.Topic]
	if e!=true {
		log.Fatalf("listTopic subsctriptions failed!\n")
		return nil,nil
	}
	var subs []string
	for _,sub :=range topic.subscri{
		subs=append(subs,sub.data.Name)
	}
	return &v1.ListTopicSubscriptionsResponse{Subscriptions:subs},nil
}

func (s *PubServer) DeleteTopic(ctx context.Context, r *v1.DeleteTopicRequest)(*emptypb.Empty,error){
	_,err:=topic_list[r.Topic]
	if err!=true {
		log.Fatalf("error when deleteTopic %s\n",r.Topic)
		return nil ,nil
	}
	delete(topic_list,r.Topic)
	return &emptypb.Empty{}, nil
}

func (s *PubServer) DetachSubscription(ctx context.Context,r *v1.DetachSubscriptionRequest)(*v1.DetachSubscriptionResponse,error){
	for _,topic := range topic_list {
		delete(topic.subscri,r.Subscription)
	}
	return &v1.DetachSubscriptionResponse{},nil
}




func main() {
	lis, err := net.Listen("tcp",":1236");
	if err!=nil {
		log.Fatalf("failed to listen: %v",err);
		return 
	}

	topic_list=make(map[string]*Topic)
	message_list=make(map[string]*v1.ReceivedMessage)
	fmt.Printf("Server booting success!\n")
	grpcServer := grpc.NewServer()
	v1.RegisterPublisherServer(grpcServer,&PubServer{ack_id:0 })
	e:=grpcServer.Serve(lis)
	if e!=nil{
		fmt.Printf("falied to serve\n")
		return 
	}
}
