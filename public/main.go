package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"strconv"

	"github.com/SkyrisBactera/govue"
	"github.com/go-humble/locstor"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"honnef.co/go/js/dom"
)

var jQuery = jquery.NewJQuery
var endpoint string
var endpoints = []string{"https://svue.psdschools.org/Service/PXPCommunication.asmx", "https://vue.d51schools.org/Service/PXPCommunication.asmx", "https://parent.ouhsd.k12.ca.us/Service/PXPCommunication.asmx", "https://d47.edupoint.com/Service/PXPCommunication.asmx", "https://afsd.edupoint.com/Service/PXPCommunication.asmx"}
var username string
var password string
var assShow bool
var changeset *govue.Changeset
var lastShown = -1
var binStore = locstor.NewDataStore(locstor.JSONEncoding)

//var grades *govue.Gradebook

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func load(key string, holder interface{}) error {
	defer timeTrack(time.Now(), "Loading "+key)
	err := binStore.Find(key, holder)
	if err != nil {
		fmt.Printf("Load Error: %s failed with\n%s\n", key, err.Error())
	}
	return err
}

func save(key string, holder interface{}) error {
	defer timeTrack(time.Now(), "Saving "+key)
	if err := binStore.Save(key, holder); err != nil {
		fmt.Printf("Save Error: %s failed with\n%s\n", key, err.Error())
		return err
	}
	return nil
}

//Really slow
func start() {
	go func() {
		defer timeTrack(time.Now(), "start")
		login()
		var oldGrades *govue.Gradebook
		if err := load("gradebook", &oldGrades); err == nil {
			go mainPage(oldGrades)
		}
		window := js.Global.Get("window")
		fmt.Println(window.Get("navigator").Get("onLine").Bool())
		if window.Get("navigator").Get("onLine").Bool() {
			tempgrades, err := govue.GetStudentGrades(username, password, endpoint)
			if err != nil {
				fmt.Println(err)
			}
			grades, err := govue.GetStudentGradesForGradingPeriod(username, password, endpoint, tempgrades.CurrentGradingPeriod.Index)
			if err != nil {
				fmt.Println(err)
			}
			if oldGrades != nil && grades != nil {
				changeSet, err := govue.CalcChangeset(oldGrades, grades)
				if err != nil {
					fmt.Println(err)
				}
				if changeSet != nil {
					go mainPage(grades)
					go afterPage(grades)
				}
			} else {
				go mainPage(grades)
				go afterPage(grades)
			}
		}
	}()
}

func main() {
	defer timeTrack(time.Now(), "main")
	js.Global.Set("svue", map[string]interface{}{
		"testAccount":     testAccount,
		"start":           start,
		"showAssignments": showAssignments,
		"toLetter":        toLetter,
	})
}

type countWG struct {
	sync.WaitGroup
	Count int // Race conditions, only for info logging.
}

func (cg *countWG) add(delta int) {
	cg.Count += delta
	cg.WaitGroup.Add(delta)
}

func (cg *countWG) done() {
	cg.Count--
	cg.WaitGroup.Done()
}

func testAccount() {
	go func() {
		defer timeTrack(time.Now(), "testAccount")
		document := dom.GetWindow().Document()
		//endpointDiv := document.GetElementByID("endpoint").(*dom.HTMLDivElement)
		username = document.GetElementByID("username").(*dom.HTMLInputElement).Value
		password = document.GetElementByID("password").(*dom.HTMLInputElement).Value
		counted := len(endpoints)
		for index := range endpoints {
			go func(index int) {
				_, err := govue.SignInStudent(username, password, endpoints[index])
				if err == nil {
					fmt.Println(endpoints[index])
					endpoint = endpoints[index]
					return
				}
				counted--
			}(index)
		}
		for endpoint == "" {
			time.Sleep(time.Millisecond)
			fmt.Print(counted)
			if counted == 0 {
				break
			}
		}
		if endpoint != "" {
			if err := save("username", username); err != nil {
				fmt.Println(err.Error())
			}
			if err := save("password", password); err != nil {
				fmt.Println(err.Error())
			}
			if err := save("endpoint", endpoint); err != nil {
				fmt.Println(err.Error())
			}
			js.Global.Get("window").Get("location").Call("replace", "index.html")
		} else {
			fmt.Println("Bad password, username, or no correct endpoint")
		}
	}()
}

func login() {
	/*
		document := dom.GetWindow().Document()
		defer timeTrack(time.Now(), "login")
		if err := load("username", &username); err != nil {
			js.Global.Get("window").Get("location").Call("replace", "/studentview/login.html")
		} else {
			document.GetElementByID("activeUser").SetInnerHTML(username)
		}
		if err := load("password", &password); err != nil {
			js.Global.Get("window").Get("location").Call("replace", "/studentview/login.html")
		}
		if err := load("endpoint", &endpoint); err != nil {
			js.Global.Get("window").Get("location").Call("replace", "/studentview/login.html")
		}
	*/
	develLogin()
}

func develLogin() {
	username = "58697"
	password = "^13371337^"
	endpoint = "https://svue.psdschools.org/Service/PXPCommunication.asmx"
}

func mainPage(grades *govue.Gradebook) {
	defer timeTrack(time.Now(), "mainPage")
	document := dom.GetWindow().Document()
	//jQuery("username").SetText(username)
	var total []float64
	fmt.Println(grades.Courses)
	var wg sync.WaitGroup
	wg.Add(len(grades.Courses))
	for index := range grades.Courses {
		go func(index int) {
			//gradeinfo := fmt.Sprintf("%s (%s)", grades.Courses[index].Teacher, grades.Courses[index].ID.Name)
			grade := grades.Courses[index].CurrentMark.RawGradeScore
			total = append(total, grade)
			wg.Done()
			lettergrade := grades.Courses[index].CurrentMark.LetterGrade
			var g dom.Element
			var bar *js.Object
			if grade != 0 {
				if document.GetElementByID(fmt.Sprintf("graph%v", index)) == nil {
					g = document.CreateElement("div")
					g.SetAttribute("id", fmt.Sprintf("graph%v", index))
					g.SetAttribute("onclick", fmt.Sprintf("svue.showAssignments(%v)", index))
					document.GetElementByID("gradegraph").AppendChild(g)
					bar = js.Global.Call("newgraph", g, "gradegraph", lettergrade)
				} else {
					g = document.GetElementByID(fmt.Sprintf("graph%v", index))
					g.SetOuterHTML("")
					g = document.CreateElement("div")
					g.SetAttribute("id", fmt.Sprintf("graph%v", index))
					g.SetAttribute("onclick", fmt.Sprintf("svue.showAssignments(%v)", index))
					document.GetElementByID("gradegraph").AppendChild(g)
					bar = js.Global.Call("newgraph", g, "gradegraph", lettergrade)
				}
				bar.Call("animate", grade/100)
			}
			//jQuery("#mainPage").Append(fmt.Sprintf("<p style='font-size: 1.5em' id='grade%v'>%s:</p><b style='color: green; font-size: 1.5em'>%s</b><hr>", index, gradeinfo, grade))
		}(index)
	}
	var bar *js.Object
	var sum float64
	var g dom.Element
	wg.Wait()
	for _, num := range total {
		sum += num
	}
	if document.GetElementByID("avggraph") == nil {
		g = document.CreateElement("div")
		g.SetID("avggraph")
		document.GetElementByID("avggradegraph").AppendChild(g)
		bar = js.Global.Call("newgraph", g, "gradegraph")
	} else {
		g = document.GetElementByID("avggraph")
		g.SetOuterHTML("")
		g = document.CreateElement("div")
		g.SetID("avggraph")
		document.GetElementByID("avggradegraph").AppendChild(g)
		bar = js.Global.Call("newgraph", g, "gradegraph")
	}
	bar.Call("animate", (sum/float64(len(total)))/100)
}

func publishChange(message string) {
	go func() {
		defer timeTrack(time.Now(), "publishChange")
		fmt.Println("New info")
		fmt.Println(message)
		document := dom.GetWindow().Document()
		g := document.CreateElement("h4")
		g.SetTextContent(message)
		br := document.CreateElement("br")
		document.GetElementByID("changedassignments").AppendChild(br)
		document.GetElementByID("changedassignments").AppendChild(g)
	}()
}

func afterPage(grades *govue.Gradebook) {
	go func() {
		defer timeTrack(time.Now(), "afterPage")
		fmt.Println("whoah")
		document := dom.GetWindow().Document()
		for index := range grades.Courses {
			go func(index int) {
				assignDiv := document.CreateElement("div")
				assignDiv.SetAttribute("style", "display:none;")
				assignDiv.SetID(fmt.Sprintf("assignments%v", index))
				for i := range grades.Courses[index].CurrentMark.Assignments {
					go func(i int) {
						assignmentP := document.CreateElement("p")
						name := grades.Courses[index].CurrentMark.Assignments[i].Name
						if !grades.Courses[index].CurrentMark.Assignments[i].Score.Graded {
							assignmentP.SetInnerHTML(fmt.Sprintf("%s: Not graded", name))
						} else if grades.Courses[index].CurrentMark.Assignments[i].ScoreType == "IB Rubric 0-8" || grades.Courses[index].CurrentMark.Assignments[i].ScoreType == "MYP Rubric Score" {
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
					}(i)
				}
				document.GetElementByID("assignments").AppendChild(assignDiv)
			}(index)
		}
		var oldGradebook govue.Gradebook
		if err := load("gradebook", &oldGradebook); err == nil {
			changeset, err = govue.CalcChangeset(&oldGradebook, grades)
			if err != nil {
				fmt.Println(err)
			}
			if changeset.CourseChanges == nil {
				jQuery("#changesornah").SetText("No changes since last time!")
			} else if changeset.CourseChanges != nil {
				fmt.Println(changeset.CourseChanges)
				document.GetElementByID("changesornah").SetAttribute("style", "display:none;")
				for index1 := range changeset.CourseChanges {
					go func(index1 int) {
						for index := range changeset.CourseChanges[index1].AssignmentChanges {
							go func(index int) {
								if changeset.CourseChanges[index1].AssignmentChanges[index].ScoreChange {
									publishChange(fmt.Sprintf("Your assignment in %s changed from %v%% to %v%%",
										changeset.CourseChanges[index1].Course.ID.Name,
										changeset.CourseChanges[index1].AssignmentChanges[index].PreviousScore,
										changeset.CourseChanges[index1].AssignmentChanges[index].NewScore))
								}
								if changeset.CourseChanges[index1].AssignmentChanges[index].PointsChange {
									publishChange(fmt.Sprintf("Your assignment in %s changed from %v%% to %v%%",
										changeset.CourseChanges[index1].Course.ID.Name,
										changeset.CourseChanges[index1].AssignmentChanges[index].PreviousPoints,
										changeset.CourseChanges[index1].AssignmentChanges[index].NewPoints))
								}
							}(index)
						}
						if changeset.CourseChanges[index1].GradeChange != nil {
							if changeset.CourseChanges[index1].GradeChange.GradeIncrease {
								publishChange(fmt.Sprintf("Your grade in %s increased from %v%% (%s) to %v%% (%s)",
									changeset.CourseChanges[index1].Course.ID.Name,
									changeset.CourseChanges[index1].GradeChange.PreviousGradePct,
									changeset.CourseChanges[index1].GradeChange.PreviousLetterGrade,
									changeset.CourseChanges[index1].GradeChange.NewGradePct,
									changeset.CourseChanges[index1].GradeChange.NewLetterGrade))
							}
							if !changeset.CourseChanges[index1].GradeChange.GradeIncrease {
								publishChange(fmt.Sprintf("Your grade in %s decreased from %v%% (%s) to %v%% (%s)",
									changeset.CourseChanges[index1].Course.ID.Name,
									changeset.CourseChanges[index1].GradeChange.PreviousGradePct,
									changeset.CourseChanges[index1].GradeChange.PreviousLetterGrade,
									changeset.CourseChanges[index1].GradeChange.NewGradePct,
									changeset.CourseChanges[index1].GradeChange.NewLetterGrade))
							}
						}
						for index := range changeset.CourseChanges[index1].AssignmentAdditions {
							go func(index int) {
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
							}(index)
						}
						for index := range changeset.CourseChanges[index1].AssignmentRemovals {
							go func(index int) {
								publishChange(fmt.Sprintf("The assignment %s was removed", changeset.CourseChanges[index1].AssignmentRemovals[index].Name))
							}(index)
						}
					}(index1)
				}
			}
		} else {
			fmt.Println(err)
		}
		if err := binStore.Delete("gradebook"); err != nil {
			// Handle err
		}
		if err := save("gradebook", grades); err != nil {
			fmt.Println("Couldn't save grades!")
			fmt.Println(err)
		}
	}()
}

func showAssignments(class string) {
	go func() {
		defer timeTrack(time.Now(), "showAssignments")
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
		js.Global.Call("scrollTo", assignmentP)
	}()
}

func toLetter(score float64) string {
	defer timeTrack(time.Now(), "toLetter")
	var grade string
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
