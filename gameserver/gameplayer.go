package main

type PlayerStatus uint8

const (
    GPS_NONE = iota
    GPS_BEGIN_SESSION
    GPS_WAITING_FOR_CG_READY
    GPS_NORMAL
    GPS_IGNORE_ALL
    GPS_AFTER_SENDING_GL_INCOMING_CONNECTION
    GPS_END_SESSION
)

type GamePlayer struct {
    Player //继承自一个Player获取网络处理方面的能力

    Creature     PlayerCreature //Ouster/Slayer/Vampire继承自Creature并实现了PlayerCreature
    PlayerStatus PlayerStatus
}
