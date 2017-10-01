package main

import (
	"fmt"

	"github.com/SkyrisBactera/govue"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"honnef.co/go/js/dom"
	//"time"
	"github.com/go-humble/locstor"
	//"reflect"
)

var jQuery = jquery.NewJQuery
var endpoint string
var endpoints []string = []string{"https://svue.psdschools.org/Service/PXPCommunication.asmx", "https://vue.d51schools.org/Service/PXPCommunication.asmx", "https://parent.ouhsd.k12.ca.us/Service/PXPCommunication.asmx", "https://d47.edupoint.com/Service/PXPCommunication.asmx", "https://afsd.edupoint.com/Service/PXPCommunication.asmx"}
var username string
var password string
var err error
var assShow bool

func start() {
	go func() {
		login()
		mainPage()
	}()
}

func main() {
	js.Global.Set("svue", map[string]interface{}{
		"testAccount":     testAccount,
		"start":           start,
		"ShowAssignments": ShowAssignments,
	})
}

func testAccount() {
	go func() { //test
		fmt.Println("Testing Account")
		document := dom.GetWindow().Document()
		//endpointDiv := document.GetElementByID("endpoint").(*dom.HTMLDivElement)
		username = document.GetElementByID("username").(*dom.HTMLInputElement).Value
		password = document.GetElementByID("password").(*dom.HTMLInputElement).Value
		for index, _ := range endpoints {
			_, err := govue.SignInStudent(username, password, endpoints[index])
			if err == nil {
				endpoint = endpoints[index]
				break
			}
		}
		if endpoint != "" {
			err = locstor.SetItem("username", username)
			if err != nil {
				go func() { js.Global.Get("window").Get("location").Call("replace", "login.html") }()
			}
			err = locstor.SetItem("password", password)
			if err != nil {
				go func() { js.Global.Get("window").Get("location").Call("replace", "login.html") }()
			}
			err = locstor.SetItem("endpoint", endpoint)
			if err != nil {
				go func() { js.Global.Get("window").Get("location").Call("replace", "login.html") }()
			}
			go func() { js.Global.Get("window").Get("location").Call("replace", "index.html") }()
		} else {
			fmt.Println("Bad password, username, or no correct endpoint")
		}
	}()
}

func login() {
	go func() {
		username, err = locstor.GetItem("username")
		if err != nil {
			go func() { js.Global.Get("window").Get("location").Call("replace", "login.html") }()
		}
		password, err = locstor.GetItem("password")
		if err != nil {
			go func() { js.Global.Get("window").Get("location").Call("replace", "login.html") }()
		}
		endpoint, err = locstor.GetItem("endpoint")
		if err != nil {
			go func() { js.Global.Get("window").Get("location").Call("replace", "login.html") }()
		}
		fmt.Println(username + password + endpoint)
	}()
}

func mainPage() {
	go func() {
		document := dom.GetWindow().Document()
		//jQuery("username").SetText(username)
		temp, _ := govue.GetStudentGrades(username, password, endpoint)
		grades, _ := govue.GetStudentGradesForGradingPeriod(username, password, endpoint, temp.CurrentGradingPeriod.Index)
		var total []float64
		for index := range grades.Courses {

			//gradeinfo := fmt.Sprintf("%s (%s)", grades.Courses[index].Teacher, grades.Courses[index].ID.Name)
			grade := grades.Courses[index].CurrentMark.RawGradeScore
			lettergrade := grades.Courses[index].CurrentMark.LetterGrade
			if grade != 0 {
				total = append(total, grade)
				fmt.Println(grade)
				g := document.CreateElement("div")
				g.SetAttribute("id", fmt.Sprintf("graph%v", index))
				g.SetAttribute("onclick", fmt.Sprintf("svue.ShowAssignments(%v)", index))
				assignDiv := document.CreateElement("div")
				assignDiv.SetAttribute("style", "display:none;")
				assignDiv.SetID(fmt.Sprintf("assignments%v", index))
				document.GetElementByID("gradegraph").AppendChild(g)
				for i := range grades.Courses[index].CurrentMark.Assignments {
					assignmentP := document.CreateElement("p")
					name := grades.Courses[index].CurrentMark.Assignments[i].Name
					score := grades.Courses[index].CurrentMark.Assignments[i].Score.Score / grades.Courses[index].CurrentMark.Assignments[i].Score.PossibleScore
					assignmentP.SetInnerHTML(fmt.Sprintf("%s: %v", name, score))
					assignDiv.AppendChild(assignmentP)
				}
				document.GetElementByID("assignments").AppendChild(assignDiv)
				fmt.Println(grades.Courses[index].ID.Name)
				bar := js.Global.Call("newgraph", g, "gradegraph", lettergrade)
				bar.Call("animate", grade/100)
			}
			//jQuery("#mainPage").Append(fmt.Sprintf("<p style='font-size: 1.5em' id='grade%v'>%s:</p><b style='color: green; font-size: 1.5em'>%s</b><hr>", index, gradeinfo, grade))
		}
		var sum float64
		for _, num := range total {
			sum += num
		}
		g := document.CreateElement("div")
		g.SetID("avggraph")
		document.GetElementByID("avggradegraph").AppendChild(g)
		bar := js.Global.Call("newgraph", g, "avggradegraph")
		bar.Call("animate", (sum/float64(len(total)))/100)
	}()
}

func ShowAssignments(class string) {
	document := dom.GetWindow().Document()
	for i := 0; i > 25; i++ {
		assignmentP := document.GetElementByID("assignments" + string(i))
		assignmentP.SetAttribute("style", "")
	}
	assignmentP := document.GetElementByID("assignments" + class)
	assignmentP.SetAttribute("style", "")
	assignments := document.GetElementByID("assignmentholder")
	assignments.SetAttribute("style", "")
}
