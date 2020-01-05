package routes

import (
	//"encoding/json"
	// "fmt"
	"github.com/go-chi/chi"
	//"net/http"
	// "github.com/likexian/whois-go"
	// "github.com/likexian/whois-parser-go"
	//"github.com/xellio/whois"
	//"net"
	//"strings"
	"apiDomainInfo/controllers"
)

//Server return the routes to acces the servers info
func Server() *chi.Mux {
	router := chi.NewRouter()
	// router.Get("/getServers", controllers.GetServers)
	router.Get("/getServerInfo/{ipAddres}", controllers.ConsultDomain)
	router.Get("/getServersRecord", controllers.GetServersRecord)
	router.Get("/holamundo", controllers.HolaMundo)
	return router
}

// //GetDomainURL2 imprime el url
// func GetDomainURL2() {
// 	result, err := whois.Whois("2a03:2880:f127:283:face:b00c:0:25de")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// }

//GetDomainURL imprime el url
// func GetDomainURL() {
// 	/*
// 		ip := net.ParseIP("8.8.8.8")
// 		res, err := whois.QueryIP(ip)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Println(strings.Join(res.Output["Registrant Country"], ","))
// 	*/
// 	results, err := whois.Whois("truora.com", "redirect1.proxy-ssl.webflow.com")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(results)
// 	result, err := whoisparser.Parse(results)
// 	if err == nil {
// 		// Print the domain status
// 		// fmt.Println(result.Registrar.DomainStatus)

// 		// Print the domain created date
// 		// fmt.Println(result.Registrar.CreatedDate)

// 		// Print the domain expiration date
// 		//fmt.Println(result.)

// 		// Print the registrant name
// 		fmt.Println(result.Registrant.Name)

// 		// Print the registrant email address
// 		// fmt.Println(result.Tech.Country)
// 		fmt.Println("Aver a ver")
// 	} else {
// 		fmt.Println(err)
// 	}
// }
