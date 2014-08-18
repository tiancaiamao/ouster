package main

type ObjectClass int

const (
    OBJECT_CLASS_CREATURE ObjectClass = iota
    OBJECT_CLASS_ITEM
    OBJECT_CLASS_OBSTACLE
    OBJECT_CLASS_EFFECT
    OBJECT_CLASS_PORTAL
)

// 每个游戏中的实体，都有一个ObjectID。客户端与服务端通过这个ID确定是什么对象
type ObjectInterface interface {
    ObjectClass() ObjectClass
    ObjectInstance() *Object
}

// 派生类，都需要继承Object对象，并实现ObjectInterface接口
type Object struct {
    ObjectID ObjectID_t
    Next     ObjectInterface
    Prev     ObjectInterface
}

func (obj Object) ObjectInstance() *Object {
    return &obj
}
