package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/jcorme/govue"
	"honnef.co/go/js/dom"
	"time"
)

var jQuery = jquery.NewJQuery
var endpoint string
var endpoints []string = []string{"https://svue.psdschools.org/Service/PXPCommunication.asmx", "https://vue.d51schools.org/Service/PXPCommunication.asmx", "https://parent.ouhsd.k12.ca.us/Service/PXPCommunication.asmx", "https://d47.edupoint.com/Service/PXPCommunication.asmx", "https://afsd.edupoint.com/Service/PXPCommunication.asmx"}
var username string
var password string

func main() {
	js.Global.Get("loginButton").Call("addEventListener", "click", func() {
		go func() {
			login()
			mainPage()
		}()
	})
	passwordText := js.Global.Get("document").Call("getElementById", "password")
	passwordText.Call("addEventListener", "keyup", handleInputKeyUp, false)
}

func handleInputKeyUp(event *js.Object) {
	go func() {
		if keycode := event.Get("keyCode").Int(); keycode == 13 {
			login()
			mainPage()
		}
	}()
}

func login() {
	document := dom.GetWindow().Document()
	endpointDiv := document.GetElementByID("endpoint").(*dom.HTMLDivElement)
	//loginDiv := document.GetElementByID("loginDiv").(*dom.HTMLDivElement)
	//feedDiv := document.GetElementByID("feedDiv").(*dom.HTMLDivElement)
	document.GetElementByID("spinner").(*dom.HTMLDivElement).Class().Toggle("invisible")
	username = document.GetElementByID("username").(*dom.HTMLInputElement).Value
	password = document.GetElementByID("password").(*dom.HTMLInputElement).Value
	for index, _ := range endpoints {
		_, err := govue.SignInStudent(username, password, endpoints[index])
		if err == nil {
			endpoint = endpoints[index]
			break
		}
	}
	document.GetElementByID("spinner").(*dom.HTMLDivElement).Class().Toggle("invisible")
	if endpoint != "" {
		jQuery("#loginDiv").FadeOut("slow", func() {
			jQuery("#feedDiv").FadeIn("slow")
		})
		time.Sleep(2 * time.Second)
		jQuery("#welcomeText").FadeOut("slow", func() {
			jQuery("#mainPage").FadeIn("slow")
		})
	} else {
		endpointDiv.Class().Toggle("invisible")
		fmt.Println("Not in a known endpoint or password is incorrect")
	}

}

func mainPage() {
	temp, _ := govue.GetStudentGrades(username, password, endpoint)
	grades, _ := govue.GetStudentGradesForGradingPeriod(username, password, endpoint, temp.CurrentGradingPeriod.Index)
	//var grade string = grades.Courses[1].CurrentMark.LetterGrade
	for index, _ := range grades.Courses {
		gradeinfo := fmt.Sprintf("%s (%s)", grades.Courses[index].Teacher, grades.Courses[index].ID.Name)
		grade := fmt.Sprintf("%v%% | (%s)", grades.Courses[index].CurrentMark.RawGradeScore, grades.Courses[index].CurrentMark.LetterGrade)
		jQuery("#mainPage").Append(fmt.Sprintf("<p style='font-size: 1.5em' id='grade%v'>%s:</p><b style='color: green; font-size: 1.5em'>%s</b><hr>", index, gradeinfo, grade))
	}
}
