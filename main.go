package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/12138mICHAEL1111/simplebank/api"
	db "github.com/12138mICHAEL1111/simplebank/db/sqlc"
	"github.com/12138mICHAEL1111/simplebank/gapi"
	"github.com/12138mICHAEL1111/simplebank/pb"
	"github.com/12138mICHAEL1111/simplebank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main(){
	config,err := util.LoadConfig("./")
	if err!=nil{
		log.Fatal("cannot load config:", err)
	}
	conn,err := sql.Open(config.DBDriver,config.DBSource)
	if err != nil {
		 log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	go runGatewayServer(config,store)
	runGrpcServer(config,store)
}

func runGrpcServer(config util.Config, store *db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot create server:",err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimplebankServer(grpcServer,server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cnnot start gRPC server")
	}
}

//使用http方法调用grpc
func runGatewayServer(config util.Config, store *db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot create server:",err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pb.RegisterSimplebankHandlerServer(ctx,grpcMux,server)

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener",err)
	}
	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener,mux)
	if err != nil {
		log.Fatal("cnnot start HTTP gateway server",err)
	}
}

func runGinServer(config util.Config, store *db.Store){
	server,err := api.NewServer(config, store)
	if err!=nil {
		log.Fatal("cannot create server",err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}