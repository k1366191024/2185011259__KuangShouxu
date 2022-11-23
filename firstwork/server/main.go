package main

import (
	"fmt"
	"net"
	v1 "firstwork/pet/v1"
	"google.golang.org/grpc"
	"google.golang.org/genproto/googleapis/type/datetime"
	"context"
	"time"
	"flag"
)

type server struct{
           v1.UnimplementedPetStoreServiceServer
}

var(
	petset map[string]v1.Pet
	id   int
)

func (s *server) GetPet(ctx context.Context, in *v1.GetPetRequest)(*v1.GetPetResponse,error){
	pet, e :=petset[in.GetPetId()]
	if e!=true{
		fmt.Printf("no aminal\n")
		return &v1.GetPetResponse{Pet:nil},nil
	}

	return &v1.GetPetResponse{Pet:&pet} ,nil
}

func(s *server) PutPet(ctx context.Context, in *v1.PutPetRequest)(*v1.PutPetResponse,error){
	ntime:=time.Now()
	uid:=fmt.Sprintf("%d",id)
	id=id+1
	npet:=v1.Pet{
		PetType: in.PetType,
		PetId: uid,
		Name: in.Name,
		CreatedAt: &datetime.DateTime{
			Year:       int32(ntime.Year()),
			Month:      int32(ntime.Month()),
			Day:        int32(ntime.Day()),
			Hours:      int32(ntime.Hour()),
			Minutes:    int32(ntime.Minute()),
			Seconds:    int32(ntime.Second()),
			Nanos:      int32(ntime.Nanosecond()),
			TimeOffset: nil,
		},
	}

	petset[uid]=npet
	return &v1.PutPetResponse{Pet:&npet},nil
}

func (s *server) DeletePet(ctx context.Context,in *v1.DeletePetRequest)(*v1.DeletePetResponse,error){
	_,e:=petset[in.GetPetId()]
	if e!=true{
		return &v1.DeletePetResponse{},nil
	}
	delete(petset,in.GetPetId())
	return &v1.DeletePetResponse{},nil
}
func main() {
	flag.Parse()
	petset=make(map[string]v1.Pet)
	id=0
	lis,e :=net.Listen("tcp",":1235")
	if e!=nil{
		fmt.Printf("failed to listes: %v",e)
		return 
	}
	fmt.Printf("Booting successfully!\nWelcome to Pet Store")
	s:=grpc.NewServer()
	v1.RegisterPetStoreServiceServer(s,&server{})
	e=s.Serve(lis)
	if e!=nil {
		fmt.Printf("failed to serve:%v",e)
		return 
	}
}
