/*
* CopyRight(C) StefanChen e-mail:2572915286@qq.com
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

package httpExService

import(
	"net/http"
	//"fmt"
	"log"
	"time"
	"common/define"
	"encoding/json"
	"strings"
)


func StartHttpService(addr string, handler define.HandlerFunc){
	log.Printf("http addr: %s.", addr)
	server := &http.Server{
		Addr:			addr,
		Handler:		handler,
		ReadTimeout:	30*time.Second,
		WriteTimeout:	30*time.Second,
		MaxHeaderBytes:	1<<16,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("http ListenAndServe: ", err)
		return
	}
}

func Get(url string, jsonst interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("htp Get, err :", err)
		return err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(jsonst)
	return err
}

func Post(url string, jsonst interface{}) error{
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("name=cjb"))
	if err != nil {
		log.Println("htp post, err :", err)
		return err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(jsonst)
	return err
}
