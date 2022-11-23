package main

import (
	"context"
	"flag"
	"log"
	"time"
	"fmt"

	v1 "Thirdwork/pet/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "127.0.0.1:1235", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewPetStoreServiceClient(conn)

	fmt.Printf("Welcome to Pet Store~ \n Please press key to operate ：\n g：get n-th pet info \n p：put new pet info \n d：delete one pet info \n q：quite the nemu\n")
	var op string 
	var none string
	_, _ = fmt.Scanln(&op)
	for op !="q"  {
		switch op{
		case "g":fmt.Printf("Please enter wanting PetId：\n")
			var pet_id string
			_, _ = fmt.Scanln(&pet_id)
			ctx, e := context.WithTimeout(context.Background(), time.Second*10)
			defer e()
			r, err := c.GetPet(ctx, &v1.GetPetRequest{PetId: pet_id})
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			log.Printf("its info:%s", r.GetPet().String())
			log.Printf("its name:%s", r.GetPet().GetName())
			log.Printf("its type:%s", r.GetPet().GetPetType())
			fmt.Printf("Press any KEY to continue\n")
			_, _ = fmt.Scanln(&none)
		case "p":fmt.Printf("Please enter its type\n 0：UNSPECIFIED\n 1：CAT\n 2：DOG\n 3：SNAKE\n 4：HAMSTER\n ")
			var pet_type int
			_, _ = fmt.Scanln(&pet_type)
			for pet_type < 0 || pet_type > 4 {
				fmt.Printf("Error! check your spelling !\n")
				fmt.Printf("Please enter its type\n 0：UNSPECIFIED\n 1：CAT\n 2：DOG\n 3：SNAKE\n 4：HAMSTER\n ")
				fmt.Scanf("%d", &pet_type)
			}
			var pet_name string
			fmt.Printf("please enter its name\n")
			_, _ = fmt.Scanln(&pet_name)
			ctx, e := context.WithTimeout(context.Background(), time.Second*10)
			defer e()
			r, err := c.PutPet(ctx, &v1.PutPetRequest{
				PetType: v1.PetType(pet_type),
				Name:    pet_name,
			})
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			log.Printf("successfully adding : %s", r.GetPet().String())
			fmt.Printf("Press any KEY to continue\n")
			_, _ = fmt.Scanln(&none)
		case "d":fmt.Printf("Please enter your deleting PetId：\n")
			var pet_id string
			_, _ = fmt.Scanln(&pet_id)
			ctx, e := context.WithTimeout(context.Background(), time.Second*10)
			defer e()
			_, err := c.DeletePet(ctx, &v1.DeletePetRequest{PetId: pet_id})
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			log.Printf("deleting: %s", pet_id)
			fmt.Printf("Press any KEY to continue\n")
			_, _ = fmt.Scanln(&none)
		case "q": 
		default: fmt.Printf("Wrong instruction! check your spelling!\n")
		}
		fmt.Printf("Welcome to Pet Store~ \n Please press key to operate ：\n g：get n-th pet info \n p：put new pet info \n d：delete one pet info \n q：quite the nemu\n")
		_, _ = fmt.Scanln(&op)
	}			
	fmt.Printf("Quiting\n")
}
