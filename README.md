# messaging-service
Messaging service for the discipline of distributed systems - Universidade Federal do Piauí - UFPI

## Job description
* Utilizando WebServices no modelo REST, desenvolva um servidor de email simplificado. Ele deve implementar, pelo menos, as seguintes funcionalidades:

- [x] Enviar mensagem - POST
- [x] Listar mensagens - GET
- [x] Apagar mensagens - DELETE
- [x] Abrir mensagem - GET
- [x] Encaminhar mensagem - PUT
- [x] Responder Mensagem - POST

* Desenvolva também um cliente que utilize o servidor através de chamadas às funcionalidades implementadas. Ao conectar, o usuário deve informar seu nome. Esta será a forma de identificação. Não é necessário preocupar-se com autenticação. As mensagens podem ser armazenadas em um simples arquivo texto. Cada mensagem deve conter, pelo menos, os seguintes campos:

- Remetente
- Destinatário
- Assunto
- Corpo

* OBS: Utilize os métodos HTTP de acordo com o que é especificado pelo modelo REST.

## Implementation details

* Este trabalho foi implementado utilizando a linguagem Go (golang) e utilizando o Windows 10. Então, siga os passos abaixo para executar este trabalho na sua máquina (tutorial feito para executar no windows 10):

- 1º baixe e instale o [Go](https://golang.org/dl/)
- 2º, observe se foi criado a variávei de ambiente "GOPATH" nas variáveis de usuário e dentro da variável "path" das variáveis de usuário foi inserido "%USERPROFILE%\go\bin" ou o caminho até o diretório bin do Go. Depois abra o prompt e execute o comando "go version" para verificar se foi instalado com sucesso.
- 3º, vá até o diretório "C:\Users\NOME_DO_USUARIO\" e procure o diretório chamado "/go". Caso não exista, crie um diretório e nomeie de "go". Dentro deste diretório criado, crie mais três diretórios, "/src", "/bin" e "/pkg".
- 4º, coloque este projeto dentro do diretório "/src".
- 5º, abra o diretório deste projeto usando o prompt e execute o comando "go run .\server.go". Depois abra o navegador e digite a URL "http://localhost:8080".

* Para mais detalhes sobre a instalação em caso de problemas, [acesse esse link](https://golang.org/doc/install).