/*
	define file data transport structure.
*/
package protocol

import (
	"io"
)

// packet implemtns pack logic
func PacketStream(writer io.Writer, reader io.Reader, size int) error {
	controlData := append([]byte(ConstHeader), IntToBytes(size)...)

	writer.Write(controlData)
    _, err := io.Copy(writer, reader)
    if err == io.EOF {
        return nil
    }
    return err
}

// Unpack implements unpack logic
func UnpackStream(buffer []byte, readerChannel chan []byte) []byte {
    length := len(buffer)
 
    var i int
    for i = 0; i < length; i = i + 1 {
        if length < i+ConstHeaderLength+ConstSaveDataLength {
            break
        }
        if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
            messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
            if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
                break
            }
            data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
            readerChannel <- data
 
            i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
        }
    }
 
    if i == length {
        return make([]byte, 0)
    }
    return buffer[i:]
}


