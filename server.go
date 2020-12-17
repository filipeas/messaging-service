package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type message struct {
	ID           string
	Remetente    string
	Destinatario string
	Assunto      string
	Corpo        string
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/login", loginHandler)                                       // UTILIZA VERBO POST
	http.HandleFunc("/newmessage", newMessageHandler)                             // UTILIZA VERBO POST
	http.HandleFunc("/sendmessage", sendMessageHandler)                           // UTILIZA VERBO GET
	http.HandleFunc("/deletemessage", deleteMessageHandler)                       // UTILIZA VERBO DELETE
	http.HandleFunc("/confirmdeletemessage", confirmDeleteMessageHandler)         // UTILIZA VERBO GET
	http.HandleFunc("/showmessage", showMessageHandler)                           // UTILIZA VERBO GET
	http.HandleFunc("/forwardmessage", forwardMessageHandler)                     // UTILIZA VERBO GET
	http.HandleFunc("/confirmforwardmessage", confirmForwardMessageHandler)       // UTILIZA VERBO PUT
	http.HandleFunc("/confirmforwardokmessage", confirmForwardOKMessageHandler)   // UTILIZA VERBO GET
	http.HandleFunc("/responsemessage", responseMessageHandler)                   // UTILIZA VERBO GET
	http.HandleFunc("/confirmresponsemessage", confirmResponseMessageHandler)     // UTILIZA VERBO GET
	http.HandleFunc("/confirmresponseokmessage", confirmResponseOKMessageHandler) // UTILIZA VERBO GET

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	name := r.FormValue("name")

	if name != "usuario1" {
		if name != "usuario2" {
			http.Error(w, "User not found.", http.StatusNotAcceptable)
			return
		}
	}

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}

func newMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/newmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	id := uuid
	name := r.FormValue("name")
	remetente := r.FormValue("remetente")
	destinatario := r.FormValue("destinatario")
	assunto := r.FormValue("assunto")
	corpo := r.FormValue("corpo")

	content, err := ioutil.ReadFile("database.txt")

	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro na leitura do arquivo: %s", err)
		return
	}

	f, err := os.Create("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em criar o arquivo: %s", err)
		return
	}
	defer f.Close()

	data := []byte(content)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
		fmt.Fprintf(w, "Erro na escrita do arquivo: %s", err2)
		return
	}

	val2 := "ID:" + id + "\nremetente:" + remetente + "\ndestinatario:" + destinatario + "\nassunto:" + assunto + "\ncorpo:" + corpo + "\n----------\n"
	data2 := []byte(val2)

	var idx int64 = int64(len(data))

	_, err3 := f.WriteAt(data2, idx)

	if err3 != nil {
		log.Fatal(err3)
	}

	// lendo mensagens do database.txt
	f2, err2 := os.Open("database.txt")
	if err2 != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f2.Close()

	scanner := bufio.NewScanner(f2)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sendmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	name := r.FormValue("name")

	// enviando dados para o sendMessage
	myvar := map[string]interface{}{"User": name}
	outputHTML(w, "static/sendMessage.html", myvar)
}

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/deletemessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		fmt.Fprintf(w, "-1")
		return
	}

	if r.Method != "DELETE" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		fmt.Fprintf(w, "-1")
		return
	}

	id := r.FormValue("id")

	// removendo mensagem do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var control = 0
	var control2 = 0
	var newMessages = ""
	for scanner.Scan() {
		if id == scanner.Text() {
			control = 1
			// fmt.Fprintf(w, scanner.Text())
			// newMessages = newMessages + scanner.Text() + "\n"
			continue
		} else {
			if scanner.Text() != "----------" {
				if control != 1 {
					// fmt.Printf(scanner.Text())
					control2 = 1
					newMessages = newMessages + scanner.Text() + "\n"
				}
			} else {
				control = 0
				if control2 != 0 {
					newMessages = newMessages + scanner.Text() + "\n"
				}
				control2 = 0
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf(newMessages)

	f2, err2 := os.Create("database.txt")
	if err2 != nil {
		log.Fatal(err2)
		fmt.Fprintf(w, "-1")
		return
	}
	defer f2.Close()

	_, err3 := f2.WriteString(newMessages)

	if err3 != nil {
		log.Fatal(err3)
		fmt.Fprintf(w, "-1")
		return
	}

	fmt.Fprintf(w, "0")
}

func confirmDeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmdeletemessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	name := r.FormValue("user")

	if name != "usuario1" {
		if name != "usuario2" {
			http.Error(w, "User not found.", http.StatusNotAcceptable)
			return
		}
	}

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}

func showMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/showmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("user")

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	message := message{}
	var achouMensagem = 0
	for _, s := range messageTemp {
		if s == id {
			achouMensagem = 1
		}

		if s != "----------" {
			if achouMensagem == 1 {
				if strings.Contains(s, "ID:") {
					message.ID = s
				}
				if strings.Contains(s, "remetente:") {
					message.Remetente = s
				}
				if strings.Contains(s, "destinatario:") {
					message.Destinatario = s
				}
				if strings.Contains(s, "assunto:") {
					message.Assunto = s
				}
				if strings.Contains(s, "corpo:") {
					message.Corpo = s
				}
			}
		} else {
			achouMensagem = 0
		}
	}
	// fmt.Println(message)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Message": message}
	outputHTML(w, "static/showMessage.html", myvar)
}

func forwardMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forwardmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("user")

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	message := message{}
	var achouMensagem = 0
	for _, s := range messageTemp {
		if s == id {
			achouMensagem = 1
		}

		if s != "----------" {
			if achouMensagem == 1 {
				if strings.Contains(s, "ID:") {
					message.ID = strings.Replace(s, "ID:", "", -1)
				}
				if strings.Contains(s, "remetente:") {
					message.Remetente = strings.Replace(s, "remetente:", "", -1)
				}
				if strings.Contains(s, "destinatario:") {
					message.Destinatario = strings.Replace(s, "destinatario:", "", -1)
				}
				if strings.Contains(s, "assunto:") {
					message.Assunto = strings.Replace(s, "assunto:", "", -1)
				}
				if strings.Contains(s, "corpo:") {
					message.Corpo = strings.Replace(s, "corpo:", "", -1)
				}
			}
		} else {
			achouMensagem = 0
		}
	}
	// fmt.Println(message)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Message": message}
	outputHTML(w, "static/forwardMessage.html", myvar)
}

func confirmForwardMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmforwardmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "PUT" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	id := uuid
	name := r.FormValue("name")
	remetente := r.FormValue("remetente")
	destinatario := r.FormValue("destinatario")
	assunto := r.FormValue("assunto")
	corpo := r.FormValue("corpo")

	content, err := ioutil.ReadFile("database.txt")

	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro na leitura do arquivo: %s", err)
		return
	}

	f, err := os.Create("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em criar o arquivo: %s", err)
		return
	}
	defer f.Close()

	data := []byte(content)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
		fmt.Fprintf(w, "Erro na escrita do arquivo: %s", err2)
		return
	}

	val2 := "ID:" + id + "\nremetente:" + remetente + "\ndestinatario:" + destinatario + "\nassunto:" + assunto + "\ncorpo:" + corpo + "\n----------\n"
	data2 := []byte(val2)

	var idx int64 = int64(len(data))

	_, err3 := f.WriteAt(data2, idx)

	if err3 != nil {
		log.Fatal(err3)
	}

	// lendo mensagens do database.txt
	f2, err2 := os.Open("database.txt")
	if err2 != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f2.Close()

	scanner := bufio.NewScanner(f2)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}

func confirmForwardOKMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmforwardokmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	name := r.FormValue("user")

	if name != "usuario1" {
		if name != "usuario2" {
			http.Error(w, "User not found.", http.StatusNotAcceptable)
			return
		}
	}

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}

func responseMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/responsemessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("user")

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	message := message{}
	var achouMensagem = 0
	for _, s := range messageTemp {
		if s == id {
			achouMensagem = 1
		}

		if s != "----------" {
			if achouMensagem == 1 {
				if strings.Contains(s, "ID:") {
					message.ID = strings.Replace(s, "ID:", "", -1)
				}
				if strings.Contains(s, "remetente:") {
					message.Remetente = strings.Replace(s, "remetente:", "", -1)
				}
				if strings.Contains(s, "destinatario:") {
					message.Destinatario = strings.Replace(s, "destinatario:", "", -1)
				}
				if strings.Contains(s, "assunto:") {
					message.Assunto = strings.Replace(s, "assunto:", "", -1)
				}
				if strings.Contains(s, "corpo:") {
					message.Corpo = strings.Replace(s, "corpo:", "", -1)
				}
			}
		} else {
			achouMensagem = 0
		}
	}
	// fmt.Println(message)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Message": message}
	outputHTML(w, "static/responseMessage.html", myvar)
}

func confirmResponseMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmresponsemessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	id := uuid
	name := r.FormValue("name")
	remetente := r.FormValue("remetente")
	destinatario := r.FormValue("destinatario")
	assunto := r.FormValue("assunto")
	corpo := r.FormValue("corpo")

	content, err := ioutil.ReadFile("database.txt")

	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro na leitura do arquivo: %s", err)
		return
	}

	f, err := os.Create("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em criar o arquivo: %s", err)
		return
	}
	defer f.Close()

	data := []byte(content)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
		fmt.Fprintf(w, "Erro na escrita do arquivo: %s", err2)
		return
	}

	val2 := "ID:" + id + "\nremetente:" + remetente + "\ndestinatario:" + destinatario + "\nassunto:" + assunto + "\ncorpo:" + corpo + "\n----------\n"
	data2 := []byte(val2)

	var idx int64 = int64(len(data))

	_, err3 := f.WriteAt(data2, idx)

	if err3 != nil {
		log.Fatal(err3)
	}

	// lendo mensagens do database.txt
	f2, err2 := os.Open("database.txt")
	if err2 != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f2.Close()

	scanner := bufio.NewScanner(f2)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}

func confirmResponseOKMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmresponseokmessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	name := r.FormValue("user")

	if name != "usuario1" {
		if name != "usuario2" {
			http.Error(w, "User not found.", http.StatusNotAcceptable)
			return
		}
	}

	// lendo mensagens do database.txt
	f, err := os.Open("database.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Erro em ler o arquivo: %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	messageTemp := []string{}
	for scanner.Scan() {
		messageTemp = append(messageTemp, scanner.Text())
	}

	messages := []message{}
	message := message{}
	for _, s := range messageTemp {
		if s != "----------" {
			if strings.Contains(s, "ID:") {
				message.ID = s
			}
			if strings.Contains(s, "remetente:") {
				message.Remetente = s
			}
			if strings.Contains(s, "destinatario:") {
				message.Destinatario = s
			}
			if strings.Contains(s, "assunto:") {
				message.Assunto = s
			}
			if strings.Contains(s, "corpo:") {
				message.Corpo = s
			}
		} else {
			messages = append(messages, message)
			message.ID = ""
			message.Remetente = ""
			message.Destinatario = ""
			message.Assunto = ""
			message.Corpo = ""
		}
	}
	// fmt.Println(messages)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// enviando dados para o dashboard
	myvar := map[string]interface{}{"User": name, "Messages": messages}
	outputHTML(w, "static/dashboard.html", myvar)
}
