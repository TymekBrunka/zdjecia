package webs

import (
	// "bufio"
	"encoding/base64"
	"encoding/json"
	// "io"
	"net/http"
	"os"
	"strings"

	// "net"
	// "fmt"
	// "encoding/json"
	"fmt"
	// "fmt"
	"errors"
	"golang.org/x/net/websocket"
	"log"
	"strconv"
	// "golang.org/x/tools/go/analysis/passes/defers"
	// "io/fs"
	"path/filepath"
)

type Img struct {
	Id   int    `json:"id"`
	Type string `json:"ft"`
	Name string `json:"name"`
	Size int    `json:"size"`
	Data string `json:"data"`
}

type FileMessage struct {
	// Type   byte  `json:"type"`;
	Sender int   `json:"nr"`
	Files  []Img `json:"files"`
	// ImgId int;
}

type Imgpair struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type WsMessage struct {
	Type    byte      `json:"type"`
	Sender  int       `json:"nr"`
	Img_ids []Imgpair `json:"ids"`
}

var Connections [2]*websocket.Conn
var cache [2]string
var Next_id int = 0
var ConLen int = 0

func HandleWebsocket(ws *websocket.Conn) {
	for {

		if ConLen < 2 || Connections[0] == ws || Connections[1] == ws {

			var msg string
			wsmsgerr := websocket.Message.Receive(ws, &msg)
			if wsmsgerr != nil {
				log.Printf("Error receiving websocket message: %s", wsmsgerr)
				Next_id = bruh(Connections[0] == ws, 0, 1).(int)
				Connections[Next_id] = nil
				ConLen--
				break
			}

			if ConLen < 2 {

				Connections[Next_id] = ws
				log.Printf("bound to id: %d", Next_id)
				Next_id = (Next_id + 1) % 2
				ConLen++
			}

			the_id := bruh(Connections[0] == ws, 0, 1).(int)

			log.Printf("(websocket) Received (id:%d): %s\n", the_id, msg)
			var data WsMessage
			errir := json.Unmarshal([]byte(msg), &data)
			if errir != nil {
				// log.Print("Websocket messsage is not json\n")
				log.Printf("\x1b[31mWebsocket message got JSON decoded wit error: %s\x1b[0m", errir)
			} else {
				log.Printf(".,decoded json message with data: \x1b[32mtype: %d, ids: %x\x1b[0m\n", data.Type, data.Img_ids)

				if data.Type == 69 {
					newmsg := fmt.Sprintf("{\"type\" : 69420, \"id\": %d, \"data\" : [ {\"nr\" : 0, \"files\" : [ ", the_id)

					files, err := os.ReadDir("./zasoby_serwera/u_0/")
					if err != nil {
						log.Printf("\x1b[31mcouldnt list directory \x1b[32mu_0\x1b[31m with error %s\x1b[0m\n", err)
					}

					for _, file := range files {
						splited := strings.Split(file.Name(), ";")
						splited[1] = strings.Replace(splited[1], "&", "/", 1)
						filePath := filepath.Join("./zasoby_serwera/u_0/", file.Name())
						newmsg = fmt.Sprintf(`%s{
                            "id"   : %s,
                            "ft"   : "%s",
                            "name" : "%s",
                            "data" : "data:%s;base64,%s"},`,
							newmsg, splited[0], splited[1], splited[2], splited[1], base64.RawStdEncoding.EncodeToString(read_file(filePath)))
					}

					newmsg = fmt.Sprintf("%s ] }, { \"nr\" : 1, \"files\": [ ", newmsg[:len(newmsg)-1])

					files, err = os.ReadDir("./zasoby_serwera/u_1/")
					if err != nil {
						log.Printf("\x1b[31mcouldnt list directory \x1b[32mu_1\x1b[31m with error %s\x1b[0m\n", err)
					}

					for _, file := range files {
						splited := strings.Split(file.Name(), ";")
						splited[1] = strings.Replace(splited[1], "&", "/", 1)
						filePath := filepath.Join("./zasoby_serwera/u_1/", file.Name())
						if err != nil {
							fmt.Println("Error opening file:", err)
							continue
						}

						newmsg = fmt.Sprintf(`%s{
                            "id"   : %s,
                            "ft"   : "%s",
                            "name" : "%s",
                            "data" : "data:%s;base64,%s"},`,
							newmsg, splited[0], splited[1], splited[2], splited[1], base64.RawStdEncoding.EncodeToString(read_file(filePath)))
					}

					newmsg = fmt.Sprintf("%s ] } ] }", newmsg[:len(newmsg)-1])

					// log.Printf(".,.,sending: \x1b[35m%s\x1b[0m\n", newmsg)

					err = websocket.Message.Send(ws, newmsg)
					if err != nil {
						log.Printf("\x1b[31mError sending websocket message: %s\x1b[0m", err)
						break
					}
				} else if data.Type == 1 {
					websocket.Message.Send(Connections[(data.Sender+1)%2], msg)
					for _, file := range data.Img_ids {
						err := remove_file(file.Id, data.Sender)
						if err != nil {
							log.Printf("\x1b[31mcouldnt remove file with id\x1b[32m%d\x1b[31m with error %s\x1b[0m\n", file.Id, err)
						}
					}
				}
			}

			// err := websocket.Message.Send(ws, msg)
			// if err != nil {
			//     log.Printf("\x1b[31mError sending websocket message: %s\x1b[0m", err)
			//     break
			// }
		} else {
			log.Println("\x1b[33mToo many connections (max 2 at a time).\x1b[0m")
			ws.Close()
			break
		}
	}
}

func decode_img_data(data_string string) ([]byte, error) {
	splits := strings.Split(data_string, ",")
	if len(splits) != 2 {
		return nil, errors.New("Invalid file base64 URI string")
	}

	data, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func bruh(cond bool, a interface{}, b interface{}) interface{} {
	if cond {
		return a
	}
	return b
}

func sanatize_ft(ft string) string {
	return strings.Replace(ft, "/", "&", 1)
}

func save_file(file Img, content []byte, sender int) error {
	folder := bruh(sender == 0, "u_0", "u_1").(string)
	filename := fmt.Sprintf("./zasoby_serwera/%s/%d;%s;%s", folder, file.Id, sanatize_ft(file.Type), file.Name)
	log.Printf(".,saving file :\x1b[32m%s\x1b[0m\n", filename)
	output, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() error {
		err := output.Close()
		return err
	}()

	// buf := make([]byte, 1024)
	if _, err := output.Write(content); err != nil {
		return err
	}
	return nil
}

func remove_file(id int, sender int) error {
	folder := bruh(sender == 0, "u_0", "u_1").(string)
	thepath := fmt.Sprintf("./zasoby_serwera/%s/", folder)
	files, err := os.ReadDir(thepath)
	if err != nil {
		// log.Printf("\x1b[31mcouldnt list directory \x1b[32m%s\x1b[31m with error %s\x1b[0m\n", folder, err)
		return err
	}

	for _, file := range files {
		id_as_string := strconv.Itoa(id)
		if strings.Split(file.Name(), ";")[0] == id_as_string {
			log.Printf(".,.,removing file \x1b[32m%s%s\x1b[0m\n", thepath, file.Name())
			err = os.Remove(fmt.Sprintf("%s%s", thepath, file.Name()))
			if err != nil {
				// log.Printf("\x1b[31mcouldnt remove file \x1b[32m%s\x1b[31m with error %s\x1b[0m\n", fmt.Sprintf("%s%s", thepath, file.Name()), err)
				return err
			}
		}
	}
	// log.Printf("removing file :\x1b[32m%s\x1b[0m\n", filename)
	return nil
}

func read_file(p string) []byte {
	file, err := os.Open(p)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return nil
	}
	fileSize := fileInfo.Size()

	// Read the file content into a byte slice
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
    return data
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
		// log.Printf("(sftc)decoded json message with data: \x1b[32msender: %d, files: %x\x1b[0m\n", msg.Sender, msg.Files)
		for _, file := range msg.Files {
			filedata, err := decode_img_data(file.Data)
			if err != nil {
				log.Printf("\x1b[31m(sftc)failed to decode file \x1b[32m%s\x1b[31m with error: %s\x1b[0m\n", file.Name, err)
			}
			// log.Printf("got file data: %s", filedata)
			err = save_file(file, filedata, msg.Sender)
			if err != nil {
				log.Printf("\x1b[31m(sftc)failed to save file \x1b[32m%s\x1b[31m with error %s\x1b[0m\n", file.Name, err)
			}
		}
		tmp, _ := json.Marshal(&msg)
		cache[(msg.Sender+1)%2] = string(tmp)
		// fmt.Printf("cache: %s\n", cache[(msg.Sender + 1) % 2])
		if Connections[(msg.Sender+1)%2] != nil {
			err := websocket.Message.Send(Connections[(msg.Sender+1)%2], "{\"type\":2137}")

			if err != nil {
				log.Printf("\x1b[31mfailed to ask other client to send request with error: %s\x1b[0m\n", err)
			} else {
				log.Printf(".,asked client %d to ask for files\n", (msg.Sender+1)%2)
			}
		}
	}
}

func UcanHasCheeseburger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	var asker int
	_, err := fmt.Fscanf(r.Body, "%d", &asker)
	if err != nil {
		log.Printf("\x1b[31mCant get asker with error: %s\x1b[0m\n", err)
	}
	log.Printf("client %d asked for files\n", asker)
	w.Write([]byte(cache[asker]))
}
