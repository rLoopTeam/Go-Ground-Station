package server

import "google.golang.org/grpc"

type SimController struct {
	conn grpc.ClientConn
}

func (client *SimController) SendCommand(){

}