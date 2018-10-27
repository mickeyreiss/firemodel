// DO NOT EDIT - Code generated by firemodel (dev).

package firemodel

import firestore "cloud.google.com/go/firestore"

type Client struct {
	Client         *firestore.Client
	TestModel      *clientTestModel
	TestTimestamps *clientTestTimestamps
}

func NewClient(client *firestore.Client) *Client {
	temp := &Client{Client: client}
	temp.TestModel = &clientTestModel{client: temp}
	temp.TestTimestamps = &clientTestTimestamps{client: temp}
	return temp
}
