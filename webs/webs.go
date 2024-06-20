package webs

import (
	// "bufio"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	// "net"
	// "fmt"
	// "encoding/json"
	"fmt"
	// "fmt"
	"errors"
	"log"

	"golang.org/x/net/websocket"
	// "golang.org/x/tools/go/analysis/passes/defers"
)

type Img struct {
    Id   int    `json:"id"`;
    Type string `json:"ft"`;
    Name string `json:"name"`;
    Size int    `json:"size"`;
    Data string `json:"data"`;
}

type FileMessage struct {
    // Type   byte  `json:"type"`;
    Sender int   `json:"nr"`;
    Files  []Img `json:"files"`;
    // ImgId int;
}

type Imgpair struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}

type WsMessage struct {
    Type    byte  `json:"type"`;
    Sender  int   `json:"nr"`;
    Img_ids []Imgpair `json:"ids"`;
}

var Connections [2]*websocket.Conn;
var cache [2]string;
var Next_id int = 0;
var ConLen int = 0;

func HandleWebsocket(ws *websocket.Conn) {
	for {

        if (ConLen < 2 || Connections[0] == ws || Connections[1] == ws) {

            var msg string
            wsmsgerr := websocket.Message.Receive(ws, &msg)
            if wsmsgerr != nil {
                log.Printf("Error receiving websocket message: %s", wsmsgerr)
                break
            }

            if (ConLen < 2) {

                Connections[Next_id] = ws
                log.Printf("bound to id: %d", Next_id)
                Next_id = (Next_id + 1) % 2
                ConLen++
            }

            the_id := bruh(Connections[0] == ws, 0, 1).(int)

            log.Printf("(websocket) Received (id:%d): %s\n", the_id, msg)
            var data WsMessage;
            errir := json.Unmarshal([]byte(msg), &data)
            if errir != nil{
                // log.Print("Websocket messsage is not json\n")
                log.Printf("\x1b[31mWebsocket message got JSON decoded wit error: %s\x1b[0m", errir)
            } else {
                log.Printf(".,decoded json message with data: \x1b[32mtype: %d, ids: %x\x1b[0m\n", data.Type, data.Img_ids)

                if (data.Type == 69){
                    newmsg := fmt.Sprintf("{\"type\" : 69420, \"id\": %d}", the_id)
                    log.Printf(".,.,sending: \x1b[35m%s\x1b[0m\n", newmsg)

                    err := websocket.Message.Send(ws, newmsg)
                    if err != nil {
                        log.Printf("\x1b[31mError sending websocket message: %s\x1b[0m", err)
                        break
                    }
                } else if (data.Type == 1){
                    websocket.Message.Send(Connections[(data.Sender + 1) % 2], msg)
                }
            }

            // err := websocket.Message.Send(ws, msg)
            // if err != nil {
            //     log.Printf("\x1b[31mError sending websocket message: %s\x1b[0m", err)
            //     break
            // }
        } else {
            log.Println("\x1b[33mToo many connections (max 2 at a time).\x1b[0m")
            ws.Close();
            break
        }
	}
}

func decode_img_data(data_string string)([]byte, error) {
    splits := strings.Split(data_string, ",")
    if (len(splits) != 2) {
        return nil, errors.New("Invalid file base64 URI string")
    }

    data, err := base64.StdEncoding.DecodeString(splits[1])
    if (err != nil) {
        return nil, err
    }

    return []byte(data), nil
}

func bruh(cond bool, a interface{}, b interface{}) (interface{}) {
    if (cond) {
        return a
    }
    return b
}

func save_file(file Img, content []byte, sender int) (error) {
    folder := bruh(sender == 0, "u_0", "u_1").(string)
    filename := fmt.Sprintf("./zasoby_serwera/%s/%s", folder, file.Name)
    log.Printf("saving file :\x1b[32m%s\x1b[0m\n", filename)
    output, err := os.Create(filename)
    if (err != nil) {
        return err
    }
    
    defer func()(error) {
        err := output.Close();
        return err;
    }()

    // buf := make([]byte, 1024)
    if _, err := output.Write(content); err != nil {
        return err
    }
    return nil
}

func SendFilesToClient(w http.ResponseWriter, r *http.Request) {
    // fmt.Printf("r: %v\n", r)
    // fmt.Printf("r: %s\n", r.Body)
    // fmt.Println("")
    var msg FileMessage
    err := json.NewDecoder(r.Body).Decode(&msg)
    if err != nil {
        log.Printf("\x1b[31m(sftc)failed to parse json request with error: %s\x1b[0m\n", err)
    } else {
        // log.Printf("recieved request (json): %+v\n", msg)
        for _, file := range msg.Files {
            filedata, err := decode_img_data(file.Data)
            if (err != nil) {
                log.Printf("\x1b[31failed to decode file \x1b[32m%s\x1b[31m with error: %s\x1b[0m\n", file.Name, err)
            }
            // log.Printf("got file data: %s", filedata)
            err = save_file(file, filedata, msg.Sender)
            if (err != nil) {
                log.Printf("\x1b[31mfailed to save file \x1b[32m%s\x1b[31m with error %s\x1b[0m\n", file.Name, err)
            }
        }
        tmp, _ := json.Marshal(&msg)
        cache[(msg.Sender + 1) % 2] = string(tmp)
        // fmt.Printf("cache: %s\n", cache[(msg.Sender + 1) % 2])
        if (Connections[(msg.Sender + 1) % 2] != nil) {
            err := websocket.Message.Send(Connections[(msg.Sender + 1) % 2], "{\"type\":2137}")
        
            if (err != nil){
                log.Printf("\x1b[31mfailed to ask other client to send request with error: %s\x1b[0m\n", err)
            } else {
                log.Printf("asked client %d to ask for files\n", (msg.Sender + 1) % 2)
            }
        }
    }
}

func UcanHasCheeseburger(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    var asker int;
    _, err := fmt.Fscanf(r.Body, "%d", &asker)
    if (err != nil) {
        log.Printf("\x1b[31mCant get asker with error: %s\x1b[0m\n", err)
    }
    log.Printf("client %d asked for files\n", asker)
    w.Write([]byte(cache[asker]))
}
