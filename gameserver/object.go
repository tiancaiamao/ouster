package main

// 每个游戏中的实体，都有一个ObjectID。客户端与服务端通过这个ID确定是什么对象
type Object struct {
    ObjectID ObjectID_t
}
