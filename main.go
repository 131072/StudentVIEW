package main

import (
	"fmt"
	"time"

	"github.com/SkyrisBactera/govue"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"honnef.co/go/js/dom"
	//"time"
	"github.com/go-humble/locstor"
	//"reflect"
	"strconv"
)

var jQuery = jquery.NewJQuery
var endpoint string
var endpoints = []string{"https://svue.psdschools.org/Service/PXPCommunication.asmx", "https://vue.d51schools.org/Service/PXPCommunication.asmx", "https://parent.ouhsd.k12.ca.us/Service/PXPCommunication.asmx", "https://d47.edupoint.com/Service/PXPCommunication.asmx", "https://afsd.edupoint.com/Service/PXPCommunication.asmx"}
var username string
var password string
var err error
var assShow bool
var changeset *govue.Changeset
var lastShown = -1

func start() {
	go func() {
		login()
		mainPage()
	}()
}

func main() {
	go func() {
		js.Global.Set("svue", map[string]interface{}{
			"testAccount":     testAccount,
			"start":           start,
			"showAssignments": showAssignments,
		})
	}()
}

func testAccount() {
	go func() { //test
		fmt.Println("Testing Account")
		document := dom.GetWindow().Document()
		//endpointDiv := document.GetElementByID("endpoint").(*dom.HTMLDivElement)
		username = document.GetElementByID("username").(*dom.HTMLInputElement).Value
		password = document.GetElementByID("password").(*dom.HTMLInputElement).Value
		for index := range endpoints {
			_, err = govue.SignInStudent(username, password, endpoints[index])
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
				g.SetAttribute("onclick", fmt.Sprintf("svue.showAssignments(%v)", index))
				assignDiv := document.CreateElement("div")
				assignDiv.SetAttribute("style", "display:none;")
				assignDiv.SetID(fmt.Sprintf("assignments%v", index))
				document.GetElementByID("gradegraph").AppendChild(g)
				for i := range grades.Courses[index].CurrentMark.Assignments {
					assignmentP := document.CreateElement("p")
					name := grades.Courses[index].CurrentMark.Assignments[i].Name
					if !grades.Courses[index].CurrentMark.Assignments[i].Score.Graded {
						assignmentP.SetInnerHTML(fmt.Sprintf("%s: Not graded", name))
					} else if grades.Courses[index].CurrentMark.Assignments[i].ScoreType == "IB Rubric 0-8" {
						proficiency := "Error"
						score := grades.Courses[index].CurrentMark.Assignments[i].Score.Score
						if score == 1 || score == 2 {
							proficiency = "Limited"
						} else if score == 3 || score == 4 {
							proficiency = "Adequate"
						} else if score == 5 || score == 6 {
							proficiency = "Proficient"
						} else if score == 7 || score == 8 {
							proficiency = "Advanced"
						}
						assignmentP.SetInnerHTML(fmt.Sprintf("%s: %v (%s)", name, grades.Courses[index].CurrentMark.Assignments[i].Score.Score, proficiency))
					} else if grades.Courses[index].CurrentMark.Assignments[i].Score.Percentage {
						assignmentP.SetInnerHTML(fmt.Sprintf("%s: %v%% (%s)", name, grades.Courses[index].CurrentMark.Assignments[i].Score.Score, toLetter(grades.Courses[index].CurrentMark.Assignments[i].Score.Score)))

					} else {
						score := grades.Courses[index].CurrentMark.Assignments[i].Score.Score / grades.Courses[index].CurrentMark.Assignments[i].Score.PossibleScore
						properscore := 100 * score
						assignmentP.SetInnerHTML(fmt.Sprintf("%s: %v%% (%s)", name, properscore, toLetter(properscore)))
					}
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
		store := locstor.NewDataStore(locstor.JSONEncoding)
		var oldGradebook govue.Gradebook
		if err := store.Find("gradebook", &oldGradebook); err == nil {
			fmt.Println("test")
			changeset, _ = govue.CalcChangeset(&oldGradebook, grades)
			if changeset.CourseChanges != nil {
				document.GetElementByID("changesornah").SetAttribute("style", "display:none;")
				for index1 := range changeset.CourseChanges {
					for index := range changeset.CourseChanges[index1].AssignmentChanges {
						if changeset.CourseChanges[index1].AssignmentChanges[index].ScoreChange {
							publishChange(fmt.Sprintf("Your assignment in %s changed from %v%% to %v%%",
								changeset.CourseChanges[index1].Course.ID,
								changeset.CourseChanges[index1].AssignmentChanges[index].PreviousScore,
								changeset.CourseChanges[index1].AssignmentChanges[index].NewScore))
						}
						if changeset.CourseChanges[index1].AssignmentChanges[index].PointsChange {
							publishChange(fmt.Sprintf("Your assignment in %s changed from %v%% to %v%%",
								changeset.CourseChanges[index1].Course.ID,
								changeset.CourseChanges[index1].AssignmentChanges[index].PreviousPoints,
								changeset.CourseChanges[index1].AssignmentChanges[index].NewPoints))
						}
					}
					if changeset.CourseChanges[index1].GradeChange.GradeIncrease {
						publishChange(fmt.Sprintf("Your grade in %s increased from %v%% (%s) to %v%% (%s)",
							changeset.CourseChanges[index1].Course.ID,
							changeset.CourseChanges[index1].GradeChange.PreviousGradePct,
							changeset.CourseChanges[index1].GradeChange.PreviousLetterGrade,
							changeset.CourseChanges[index1].GradeChange.NewGradePct,
							changeset.CourseChanges[index1].GradeChange.NewLetterGrade))
					}
					if !changeset.CourseChanges[index1].GradeChange.GradeIncrease {
						publishChange(fmt.Sprintf("Your grade in %s decreased from %v%% (%s) to %v%% (%s)",
							changeset.CourseChanges[index1].Course.ID,
							changeset.CourseChanges[index1].GradeChange.PreviousGradePct,
							changeset.CourseChanges[index1].GradeChange.PreviousLetterGrade,
							changeset.CourseChanges[index1].GradeChange.NewGradePct,
							changeset.CourseChanges[index1].GradeChange.NewLetterGrade))
					}
					for index := range changeset.CourseChanges[index1].AssignmentAdditions {
						message := fmt.Sprintf("A new assignment called %s was added. It was assigned on %v.",
							changeset.CourseChanges[index1].AssignmentAdditions[index].Name,
							changeset.CourseChanges[index1].AssignmentAdditions[index].Date.Time)
						if !changeset.CourseChanges[index1].AssignmentAdditions[index].DueDate.IsZero() {
							if changeset.CourseChanges[index1].AssignmentAdditions[index].DueDate.Time.After(time.Now()) {
								message = message + fmt.Sprintf(" It is due on %v.", changeset.CourseChanges[index1].AssignmentAdditions[index].DueDate.Time)
							} else {
								message = message + fmt.Sprintf(" It was due on %v.", changeset.CourseChanges[index1].AssignmentAdditions[index].DueDate.Time)
							}
						}
						if changeset.CourseChanges[index1].AssignmentAdditions[index].Score.Graded {
							message = message + fmt.Sprintf(" You got a %v%%.", 100*(changeset.CourseChanges[index1].AssignmentAdditions[index].Points.Points/changeset.CourseChanges[index1].AssignmentAdditions[index].Points.PossiblePoints))
						}
						if changeset.CourseChanges[index1].AssignmentAdditions[index].Points.Points != 0 {
							message = message + fmt.Sprintf(" It is worth %v points", changeset.CourseChanges[index1].AssignmentAdditions[index].Points.Points)
						}
						if changeset.CourseChanges[index1].AssignmentAdditions[index].Notes != "" {
							message = message + fmt.Sprintf(" %s added a note, \"%s\".", changeset.CourseChanges[index1].Course.Teacher, changeset.CourseChanges[index1].AssignmentAdditions[index].Notes)
						}
						publishChange(message)
					}
					for index := range changeset.CourseChanges[index1].AssignmentRemovals {
						publishChange(fmt.Sprintf("The assignment %s was removed", changeset.CourseChanges[index1].AssignmentRemovals[index].Name))
					}
				}
			}
		} else {
			fmt.Println(err)
		}
		if err := store.Delete("gradebook"); err != nil {
			// Handle err
		}
		if err := store.Save("gradebook", grades); err != nil {
			fmt.Println("Couldn't save grades!")
			fmt.Println(err)
		}
	}()
}

func publishChange(message string) {
	go func() {
		document := dom.GetWindow().Document()
		g := document.CreateElement("h4")
		br := document.CreateElement("br")
		document.GetElementByID("changedassignments").AppendChild(br)
		document.GetElementByID("changedassignments").AppendChild(g)
	}()
}

func showAssignments(class string) {
	go func() {
		fmt.Println("asked to show assignments")
		document := dom.GetWindow().Document()
		if lastShown != -1 {
			assignmentP := document.GetElementByID("assignments" + strconv.Itoa(lastShown))
			assignmentP.SetAttribute("style", "display:none;")
		}
		assignmentP := document.GetElementByID("assignments" + class)
		assignmentP.SetAttribute("style", "")
		assignments := document.GetElementByID("assignmentholder")
		assignments.SetAttribute("style", "")
		lastShown, _ = strconv.Atoi(class)
	}()
}

func toLetter(score float64) string {
	grade := "Couldn't get letter grade"
	fmt.Println("Processing: " + strconv.FormatFloat(score, 'f', 6, 64))
	if score >= 100 {
		grade = "A+"
	} else if score >= 93 && score < 100 {
		grade = "A"
	} else if score <= 92.9 && score >= 89 {
		grade = "A-"
	} else if score <= 88.9 && score >= 87 {
		grade = "B+"
	} else if score <= 86.9 && score >= 83 {
		grade = "B"
	} else if score <= 82.9 && score >= 79 {
		grade = "B-"
	} else if score <= 78.9 && score >= 77 {
		grade = "C+"
	} else if score <= 76.9 && score >= 73 {
		grade = "C"
	} else if score <= 72.9 && score >= 69 {
		grade = "C-"
	} else if score <= 68.9 && score >= 67 {
		grade = "D+"
	} else if score <= 66.9 && score >= 60 {
		grade = "D"
	} else {
		grade = "F"
	}
	return grade
}
