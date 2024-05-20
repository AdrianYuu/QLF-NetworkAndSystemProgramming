package main

import (
	"fmt"
	"net"
	"net/http"
)

func middleware(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request){
	_, _ = fmt.Fprintf(w, "\nSuccessfully update!\n")
}

func addHandler(w http.ResponseWriter, r *http.Request){
	_, _ = fmt.Fprintf(w, "\nSuccessfully add!\n")
}

func informationHandler(w http.ResponseWriter, r *http.Request){
	_, _ = fmt.Fprintf(w, "\nCountry: Indonesia\nCity: Palembang\nPhone: 081234567890\n")
}

func main(){
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/store", middleware(http.MethodGet, informationHandler))
	serveMux.HandleFunc("/store/add", middleware(http.MethodPost, addHandler))
	serveMux.HandleFunc("/store/update-stock", middleware(http.MethodPut, updateHandler))
	
	httpServer := &http.Server{
		Addr: "localhost:8888",
		Handler: serveMux,
	}

	listener, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	err = httpServer.Serve(listener)

	if err != http.ErrServerClosed {
		fmt.Println(err)
		return
	}
}