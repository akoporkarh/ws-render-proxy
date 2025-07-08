package main

import (
    "io"
    "log"
    "net/http"
    "net/url"

    "github.com/gorilla/websocket"
)

var backend = "ws://45.91.201.220:10000/ws" // آدرس سرور اصلی Xray

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
    // اتصال ورودی
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrade:", err)
        return
    }
    defer c.Close()

    u, _ := url.Parse(backend)
    d, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Println("dial:", err)
        return
    }
    defer d.Close()

    // دوطرفه کردن اتصال
    go io.Copy(c.UnderlyingConn(), d.UnderlyingConn())
    io.Copy(d.UnderlyingConn(), c.UnderlyingConn())
}

func main() {
    http.HandleFunc("/", handler)
    log.Println("Listening on :10000")
    log.Fatal(http.ListenAndServe(":10000", nil))
}