package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	//conexion database
	_ "github.com/lib/pq"
	"github.com/likexian/whois-go"
)

//Server is the struct of a server
type Server struct {
	Address  string `json:"addres"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

//HostMain return server record
type HostMain struct {
	Items []dominio
}

type dominio struct {
	Dominio string
	Info    Host
}

//Host is the struct of a host
type Host struct {
	Endpoints        []endpointsInfo
	ServersChanged   bool
	SslGrade         string
	PreviousSslGrade string
	Logo             string
	Title            string
	IsDown           bool
}

type endpointsInfo struct {
	IPAddress string
	Grade     string
	Country   string
	Owner     string
}

func HolaMundo(res http.ResponseWriter, req *http.Request){
	fmt.Println("Hola papa")
}

//GetServersRecord return record servers
func GetServersRecord(res http.ResponseWriter, req *http.Request) {
	var recordServers HostMain
	if canConnect() {
		recordServers = getServersRecord()
	}
	json.NewEncoder(res).Encode(recordServers)
}

//ConsultDomain params a domain return server data || Validaciones pendientes, cambiar Fatalln por fmt.println
func ConsultDomain(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	ipAdd := chi.URLParam(req, "ipAddres")
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + ipAdd) //Validar
	if err != nil {
		log.Println("err consult domain ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err read body ", err)
	}

	var serverInfo Host
	lessGrade := "A+"
	json.Unmarshal(body, &serverInfo)
	//serverInfo.Endpoints[0].Grade = "F" //Datos para pruebas
	//serverInfo.Endpoints[1].Grade = "C"
	//asignando pais, propietario y menor grado ssl
	for i := 0; i < len(serverInfo.Endpoints); i++ {
		p, o := getOwnerServer(serverInfo.Endpoints[i].IPAddress)
		serverInfo.Endpoints[i].Country = p
		serverInfo.Endpoints[i].Owner = o
		fmt.Println(lessGrade)
		if getNumberGrade(serverInfo.Endpoints[i].Grade) < getNumberGrade(lessGrade) {
			lessGrade = serverInfo.Endpoints[i].Grade
			serverInfo.SslGrade = lessGrade
		}
	}
	//Obtenemos el logo y el titulo de la pagina, en caso de no obtenerse, se especifica
	serverInfo.Title, serverInfo.Logo = getLogoAndTitle(ipAdd) //pasar un domino invalido para probar el manejo de errores
	//Obtenemos el grado menor anterior, en caso de ser un nuevo registro, se retorna ""

	//---Base de Datos----
/* 	if canConnect() {
		serverInfo.PreviousSslGrade = getSslGrade(ipAdd)

		if !isDomainInDataBase(ipAdd) {
			saveDomainDataBase(ipAdd, serverInfo)
		} else {
			updateDomain(ipAdd, serverInfo)
		}

	} */
	//---/Base de Datos----

	if _, err := http.Get("https://" + ipAdd); err != nil {
		log.Println("err get logo and title ", err)
		serverInfo.IsDown = true
	}
	//fmt.Println("SI LLEGO AQUI")
	// json.NewEncoder(res).Encode(serverInfo)
	response, _ := json.Marshal(serverInfo)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(202)
	res.Write(response)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//revisar country con la otra libreria
func getOwnerServer(ip string) (country, owner string) {
	result, err := whois.Whois(ip)
	if err != nil {
		fmt.Println(err)
	}

	country, owner = "not found", "not found"

	values := strings.Split(result, "\n")
	for _, text := range values {
		if strings.Contains(text, "Country") {
			country = strings.ReplaceAll(strings.ReplaceAll(text, "  ", ""), "Country:", "")
		} else if strings.Contains(text, "country") {
			country = strings.ReplaceAll(strings.ReplaceAll(text, "  ", ""), "country:", "")
		}
		if strings.Contains(text, "OrgName:") {
			owner = strings.ReplaceAll(strings.ReplaceAll(text, "  ", ""), "OrgName:", "")
		} else if strings.Contains(text, "org-name:") {
			owner = strings.ReplaceAll(strings.ReplaceAll(text, "  ", ""), "org-name:", "")
		}
		if country != "not found" && owner != "not found" {
			return
		}
	}
	return
}

func getNumberGrade(grade string) int {
	ssLabs := [8]string{"F", "E", "D", "C", "B", "A", "A+", ""}
	for index, value := range ssLabs {
		if value == grade {
			return index
		}
	}
	return -1
}

//Validaciones pendientes y codigo comentado por borrar
func getLogoAndTitle(domain string) (title, logo string) {
	resp, err := http.Get("https://" + domain) //Validar
	if err != nil {
		log.Println("err get logo and title ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err read body logo and title ", err)
	}
	html := string(body)
	posI := strings.Index(html, "<title")
	title = "title not found"
	if posI != -1 {
		htmlCopy := html[posI+6 : len(html)]
		nextM := strings.Index(htmlCopy, ">")
		htmlCopy = htmlCopy[nextM+1 : len(htmlCopy)]
		nextM = strings.Index(htmlCopy, "</")
		title = htmlCopy[:nextM]
	}
	posI = strings.Index(html, "shortcut icon")
	logo = "logo not found"
	/*
		La estrategia para obtener le icono, es obtener la etiqueta <link> donde esta almacenado el link del logo,
		se obtine mendiante la busqueda de la palabra "shortcut icon", cuando esta se haya, se sabe que esta dentro
		de la etiqueta <link>. hmtl es un string que contiene todo el HTML de la pagina, ubico el index de "shortcut icon",
		busco la primer etiqueta de cierre ">" despues de "shortcut icon", luego invierto la cadena y hayo la etiqueta de
		apertuta "<", en este momento ya tengo todo el contenido de <link>, de ahi extraigo el href
	*/
	if posI != -1 {
		posI := strings.Index(html, "shortcut icon")
		htmlCopy := html[posI:len(html)]
		posM := strings.Index(htmlCopy, ">")
		htmlCopy = htmlCopy[:posM+1]
		htmlCopy2 := reverse(html)
		posR := strings.Index(htmlCopy2, reverse("shortcut icon"))
		htmlCopy2 = htmlCopy2[posR:len(html)]
		posRM := strings.Index(htmlCopy2, "<")
		htmlCopy2 = htmlCopy2[2 : posRM+1] //En lugar de dos, debrias ser len("shortcut icon")
		htmlCopy2 = reverse(htmlCopy2)
		htmlCopy2 += htmlCopy
		posIni := strings.Index(htmlCopy2, "href")
		htmlCopy2 = htmlCopy2[posIni+6 : len(htmlCopy2)]
		logo = htmlCopy2[:strings.Index(htmlCopy2, "\"")]
	}
	return
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func canConnect() bool {
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos", err)
		return false
	}
	defer db.Close()

	if _, err := db.Prepare(""); err != nil {
		log.Println("Fallo al conectar basedatos filtro 2", err)
		return false
	}

	fmt.Println("Si trato de Conectarse")
	return true
}

func getServersRecord() HostMain {
	var recordHost HostMain

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos gSR", err)
	}
	defer db.Close()

	stmtServersRecord, err := db.Prepare(`SELECT domain,servers_changed,ssl_grade,previous_ssl_grade,logo,title,is_down FROM domain_info`)
	if err != nil {
		log.Println("Err select prepare statement domain_info ", err)
	}
	defer stmtServersRecord.Close()

	rows, err := stmtServersRecord.Query()
	if err != nil {
		log.Println("Err consult domain_info ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var domains dominio
		if err := rows.Scan(&domains.Dominio, &domains.Info.ServersChanged, &domains.Info.SslGrade, &domains.Info.PreviousSslGrade, &domains.Info.Logo, &domains.Info.Title, &domains.Info.IsDown); err != nil {
			log.Println("err to get domain_info ", err)
		}
		domains.Info.Endpoints = getEndpoint(domains.Dominio)
		recordHost.Items = append(recordHost.Items, domains)
	}

	return recordHost
}

func getEndpoint(domain string) []endpointsInfo {
	var endPs []endpointsInfo
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos gSR", err)
	}
	defer db.Close()

	stmtEnpoints, err := db.Prepare(`SELECT address, ssl_grade, country, owner FROM servers_info WHERE domain = $1`)
	if err != nil {
		log.Println("Err select prepare statement servers_info ", err)
	}
	defer stmtEnpoints.Close()

	rows, err := stmtEnpoints.Query(domain)
	if err != nil {
		log.Println("Err consult domain_info ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var endP endpointsInfo
		if err := rows.Scan(&endP.IPAddress, &endP.Grade, &endP.Country, &endP.Owner); err != nil {
			log.Println("err to get endpoint ", err)
		}
		endPs = append(endPs, endP)
	}

	return endPs
}

func updateDomain(domain string, domainInfo Host) {
	//conexion a la base de datos
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos uD", err)
	}
	defer db.Close()

	//create prepare statement to update table domain_info
	stmtInsertDomainInfo, err := db.Prepare(`UPDATE domain_info SET domain = $1, servers_changed = $2, ssl_grade = $3, previous_ssl_grade = $4, logo = $5, title = $6, is_down = $7 WHERE domain = $8;`)
	if err != nil {
		log.Println("Err create prepare statement to update domain_info ", err)
	}
	defer stmtInsertDomainInfo.Close()

	stmtInsertServerInfo, err := db.Prepare(`UPDATE servers_info SET address = $1, ssl_grade = $2, country = $3, owner = $4 WHERE address = $1;`)
	if err != nil {
		log.Println("Err create prepare statement to update servers_info ", err)
	}
	defer stmtInsertServerInfo.Close()

	if _, err := stmtInsertDomainInfo.Exec(domain, domainInfo.ServersChanged, domainInfo.SslGrade, domainInfo.PreviousSslGrade, domainInfo.Logo, domainInfo.Title, domainInfo.IsDown, domain); err != nil {
		log.Println("Err update data domain_info ", err)
	}

	//update servers info on database
	for _, value := range domainInfo.Endpoints {
		if _, err := stmtInsertServerInfo.Exec(value.IPAddress, value.Grade, value.Country, value.Owner); err != nil {
			log.Println("Err update data servers_info ", err)
		}
	}

}

func saveDomainDataBase(domain string, domainInfo Host) {
	//conexion a la base de datos
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos sDB", err)
	}
	defer db.Close()

	//create prepare statement to insert data in table domain_info
	stmtInsertDomainInfo, err := db.Prepare(`INSERT INTO domain_info (domain,servers_changed,ssl_grade,previous_ssl_grade,logo,title,is_down) VALUES ($1, $2, $3, $4, $5, $6, $7);`)
	if err != nil {
		log.Println("Err create prepare statement to insert domain_info ", err)
	}
	defer stmtInsertDomainInfo.Close()

	//create prepare statement to insert data in table servers_info
	stmtInsertServerInfo, err := db.Prepare(`INSERT INTO servers_info (address,ssl_grade,country,owner,domain) VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		log.Println("Err create prepare statement to insert servers_info ", err)
	}
	defer stmtInsertServerInfo.Close()

	//Insert domain info on database
	if _, err := stmtInsertDomainInfo.Exec(domain, domainInfo.ServersChanged, domainInfo.SslGrade, domainInfo.PreviousSslGrade, domainInfo.Logo, domainInfo.Title, domainInfo.IsDown); err != nil {
		log.Println("Err insert data domain_info ", err)
	}

	//Insert servers info on database
	for _, value := range domainInfo.Endpoints {
		if _, err := stmtInsertServerInfo.Exec(value.IPAddress, value.Grade, value.Country, value.Owner, domain); err != nil {
			log.Println("Err insert data servers_info ", err)
		}
	}
}

func isDomainInDataBase(domain string) bool {
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos iDIDB", err)
	}
	defer db.Close()

	stmtIsDomain, err := db.Prepare(`SELECT 1 FROM domain_info WHERE domain = $1`)
	if err != nil {
		log.Println("Err create prepare statement to consult domain ", err)
	}
	defer stmtIsDomain.Close()

	rows, err := stmtIsDomain.Query(domain)
	if err != nil {
		log.Println("Err consult domain ", err)
	}
	defer rows.Close()

	return rows.Next()
}

func getSslGrade(domain string) (sslGrade string) {
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/testtech?sslmode=disable")
	if err != nil {
		log.Println("Fallo al conectar basedatos gSG", err)
	}
	defer db.Close()

	stmtIsSslGrade, err := db.Prepare(`SELECT ssl_grade FROM domain_info WHERE domain = $1`)
	if err != nil {
		log.Println("Err create prepare statement to consult ssl_grade ", err)
	}
	defer stmtIsSslGrade.Close()

	rows, err := stmtIsSslGrade.Query(domain)
	if err != nil {
		log.Println("Err consult ssl_grade ", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&sslGrade); err != nil {
			log.Println("err to get ssl_grade")
		}
	}
	return
}

//GetServers return all servers
// func GetServers(res http.ResponseWriter, req *http.Request) {
// 	var server Server
// 	server.Address = "direccion"
// 	server.SslGrade = "ssl_detalles"
// 	server.Country = "Country"
// 	server.Owner = "ownwe"
// 	json.NewEncoder(res).Encode(server)
// }
