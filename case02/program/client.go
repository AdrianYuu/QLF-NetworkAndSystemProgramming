package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func addMenu(){
	var mangaName string
	var mangaGenre string
	var mangaStock int

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Input manga name: ")
		mangaName, _ = reader.ReadString('\n')
		mangaName = strings.TrimSpace(mangaName)

		if mangaName != "" {
			break
		} else {
			fmt.Println("Manga name is required!")
		}
	}

	for {
		fmt.Print("Input manga genre [Fantasy | Action | Romance (case sensitive)]: ")
		mangaGenre, _ = reader.ReadString('\n')
		mangaGenre = strings.TrimSpace(mangaGenre)

		if mangaGenre == "Fantasy" || mangaGenre == "Action" || mangaGenre == "Romance" {
			break
		} else {
			fmt.Println("Manga genre must be 'Fantasy' or 'Action' or 'Romance'!")
		}
	}

	for {
		fmt.Print("Input manga stock [must greater than 0]: ")
		fmt.Scanf("%d\n", &mangaStock)

		if mangaStock > 0 {
			break
		} else {
			fmt.Println("Manga stock must be greater than 0!")
		}
	}

	reqBody := new(bytes.Buffer)
	w := multipart.NewWriter(reqBody)

	nameField, err := w.CreateFormField("Name")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = nameField.Write([]byte(mangaName))
	if err != nil {
		fmt.Println(err)
		return
	}

	genreField, err := w.CreateFormField("Genre")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = genreField.Write([]byte(mangaGenre))
	if err != nil {
		fmt.Println(err)
		return
	}

	stockField, err := w.CreateFormField("Stock")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = stockField.Write([]byte(string(mangaStock)))
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Open("file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fileField, err := w.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(fileField, file)
	if err != nil {
		fmt.Println(err)
		return
	} 

	w.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 8 * time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8888/store/add", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout(){
			fmt.Println("Request timed out!")
		} else {
			fmt.Println(err)
		}
		return
	}

	defer func(){
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func updateStock(){
	var updatedNewStock int

	for {
		fmt.Print("Input new stock [must be greater 0]: ")
		fmt.Scanf("%d\n", &updatedNewStock)

		if updatedNewStock > 0 {
			break
		} else {
			fmt.Println("Manga stock must be greater than 0!")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8 * time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "http://localhost:8888/store/update-stock", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout(){
			fmt.Println("Request timed out!")
		} else {
			fmt.Println(err)
		}
		return
	}

	defer func(){
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func updateMenu(){
	var choice int

	for {
		fmt.Println("Manga List: ")
		fmt.Println("=====================================")
		fmt.Println("| No. | Name               | Stocks |")
		fmt.Println("=====================================")
		fmt.Println("| 1.  | Frieren            |   60   |")
		fmt.Println("| 2.  | Attack On Titan    |   40   |")
		fmt.Println("| 3.  | Jujutsu Kaisen     |   25   |")
		fmt.Println("| 4.  | Violet Evergarden  |   10   |")
		fmt.Println("| 5.  | Sword Art Online   |   35   |")
		fmt.Println("=====================================")
		fmt.Println("[0] to Exit.")
		fmt.Print(">> ")
		
		_, err := fmt.Scanf("%d\n", &choice)

		if err != nil {
			fmt.Println("Input must be a number [integer].")
		}

		if choice < 0 || choice > 5 {
			fmt.Println("Invalid input.")
		} else if choice == 0 {
			fmt.Println("Thank you.")
			break;
		} else {
			updateStock();
		}
	}
}

func informationMenu(){
	ctx, cancel := context.WithTimeout(context.Background(), 8 * time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8888/store", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout(){
			fmt.Println("Request timed out!")
		} else {
			fmt.Println(err)
		}
		return
	}

	defer func(){
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func main(){
	var choice int

	for {
		fmt.Println("AY Manga Store")
		fmt.Println("==============")
		fmt.Println("1. Add product")
		fmt.Println("2. Update product")
		fmt.Println("3. Store information")
		fmt.Println("0. Exit")
		fmt.Print(">> ")
		_, err := fmt.Scanf("%d\n", &choice)

		if err != nil {
			fmt.Println("Input must be a number [integer].")
		}

		if choice == 1 {
			addMenu()
		} else if choice == 2 {
			updateMenu()
		} else if choice == 3 {
			informationMenu()
		} else if choice == 0 {
			break
		} else if choice < 0 || choice > 3 {
			fmt.Println("Invalid input.")
		}
	}
}