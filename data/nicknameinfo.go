package data

import (
    "encoding/binary"
    "errors"
    "io"
)

type NicknameInfo struct {
    NicknameID    uint16
    NicknameType  uint8
    Nickname      string
    NicknameIndex uint16
}

func (info *NicknameInfo) Size() uint32 {
    return 3
}

func (info *NicknameInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, info.NicknameID)
    binary.Write(writer, binary.LittleEndian, uint8(0))
    return nil
}

func (info *NicknameInfo) Read(reader io.Reader) error {
    binary.Read(reader, binary.LittleEndian, &info.NicknameID)
    binary.Read(reader, binary.LittleEndian, &info.NicknameType)

    switch info.NicknameType {
    case 0:
        break
    default:
        return errors.New("not implement!")
    }
    return nil
}
