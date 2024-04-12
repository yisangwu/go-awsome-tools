package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go-awsome-tools/model"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
    h bool
    db_host    string
    db_port    int
    db_user    string
    db_password string
    database string
    table string
    file_name   string
)

func init() {
    flag.BoolVar(&h, "h", false, "show help")
    flag.StringVar(&db_host, "db_host", "localhost", "databse host")
    flag.IntVar(&db_port, "db_port", 3306, "databse port")
    flag.StringVar(&db_user, "db_user", "root", "databse user")
    flag.StringVar(&db_password, "db_password", "123456", "databse db_password")
    flag.StringVar(&database, "database", "test", "databse name")
    flag.StringVar(&table, "table", "test", "table name")
    flag.StringVar(&file_name, "file_name", "xx.csv", "In the current directory, file name")

    flag.Usage = usage
}

func usage() {
    fmt.Fprintf(os.Stderr, `
Analyze question bank files and write them to the database
Usage: QuizeCsvImport -db_host aaa / --db_host aaa / -db_host=aaa / --database=aaaa

Options:
`)
    flag.PrintDefaults()
}

func GenModel(con *gorm.DB, table string) {
    g := gen.NewGenerator(gen.Config{
        OutPath: "models",
        OutFile: "go",
        Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface, 
        FieldWithTypeTag: true,
        FieldWithIndexTag: false,
      })

    g.UseDB(con)
    g.GenerateModel(table)
    g.Execute()
}


func ParseQuizeFile(filename string) []*model.QuestionBank{

    quize:= []*model.QuestionBank{}

    f_handle, err := os.Open(filename)
    if err!=nil{
        fmt.Printf("open file failed, filename:%s, err:%+v", filename, err)
        return quize
    }
    defer f_handle.Close()

    fileScanner := bufio.NewScanner(f_handle)
    for fileScanner.Scan() {
        line := fileScanner.Text()
        if line == "" {
            continue
        }
        str_arr := strings.Split(string(line), `	`)
        lang_id,_ := strconv.Atoi(str_arr[4])
        correct_answer,_ := strconv.Atoi(str_arr[3])
        quize = append(quize, &model.QuestionBank{
            LangID: int32(lang_id),
            Question: str_arr[0],
            Option1: str_arr[1],
            Option2: str_arr[2],
            CorrectAnswer: int32(correct_answer),
        })
    }
    return quize
}

type JsonDataStruct struct {
    LangID        int    `json:"lang_id"`
    Question      string `json:"question"`
    Option1       string `json:"option_1"`
    Option2       string `json:"option_2"`
    CorrectAnswer int    `json:"correct_answer"`
}


func ParseQuizeFileToJson(filename string)  {
    
    f_handle, err := os.Open(filename)
    if err!=nil{
        fmt.Printf("open file failed, filename:%s, err:%+v", filename, err)
        return
    }
    defer f_handle.Close()

    fileScanner := bufio.NewScanner(f_handle)
    quize_arr := []JsonDataStruct{}
    for fileScanner.Scan() {
        line := fileScanner.Text()
        if line == "" {
            continue
        }
        str_arr := strings.Split(string(line), `	`)
        lang_id,_ := strconv.Atoi(str_arr[4])
        correct_answer,_ := strconv.Atoi(str_arr[3])
        quize_arr = append(quize_arr, JsonDataStruct{
            LangID: int(lang_id),
            Question: str_arr[0],
            Option1: str_arr[1],
            Option2: str_arr[2],
            CorrectAnswer: int(correct_answer),
        })
    }
    quize_json, err:= json.Marshal(quize_arr)
    if err!=nil{
        fmt.Printf("ParseQuizeFileToJson json marshal failed, err:%+v", err)
        return
    }
    if err := os.WriteFile("questions.json", quize_json, 0666); err != nil {
        fmt.Printf("ParseQuizeFileToJson WriteFile failed, err:%+v", err)
        return
    }
}


func SelectFromTableMarshalToFile(conn *gorm.DB){
    
    question_bank := []model.QuestionBank{}
    allData := conn.Find(&question_bank)
    if allData.Error != nil{
        fmt.Printf("select failed, err:%+v", allData.Error)
        return
    }
    xx := `[
["题库信息总表"],
["Id", "languageID", "questionKey", "option1Key", "option2Key", "answerKey", "answerID"],
["int", "int", "string", "string", "string", "string", "int"],
["题目id","所属语言id", "问题", "选项1", "选项2", "答案", "答案id"],`+"\n"
    quize_arr := []JsonDataStruct{}
    len_bank := len(question_bank)
    for i, item := range question_bank{
        fmt.Println(item.LangID)
        quize_arr = append(quize_arr, JsonDataStruct{
            LangID: int(item.LangID),
            Question: item.Question,
            Option1: item.Option1,
            Option2: item.Option2,
            CorrectAnswer: int(item.CorrectAnswer),
        } )

        jsonStr, _ := json.Marshal([]string{
            strconv.Itoa(int(item.ID)), 
            strconv.Itoa(int(item.LangID)), 
            item.Question, 
            item.Option1, 
            item.Option2, 
            strconv.Itoa(int(item.CorrectAnswer)),
        })
        if (len_bank -i ) >1 {
            xx += string(jsonStr) + ",\n"
        }else{
            xx += string(jsonStr) + "\n"
        }
    }
    xxxx, _ := json.Marshal(quize_arr)
    fmt.Println(string(xxxx))
   
    xx = strings.TrimRight(xx, ",")
    xx += "]"
    fmt.Println(string(xx))

    if err := os.WriteFile("questions.json", []byte(xx), 0666); err != nil {
        fmt.Printf("SelectFromTableMarshalToFile WriteFile failed, err:%+v", err)
        return
    }

}


// go run .\QuizeCsvImport.go -db_host="192.168.1.222" -db_user=admin -db_password=admin -database=quiz_master -table=question_bank -file_name="quize.txt"
// xx.txt ： Qual é a dança da paixão nacional do Brasil?	tango	Samba	2	2
func main() {
    flag.Parse()
    if h {
        flag.Usage()
    }
    dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_password, db_host, db_port, database)
    fmt.Printf("dsn:%s", dsn)

    conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil{
        fmt.Printf("connect database failed, err:%+v", err)
        os.Exit(1)
    }
    // GenModel(conn, table)

    // bluckModels := ParseQuizeFile(file_name)
    // if len(bluckModels)==0 {
    //     fmt.Printf("ParseQuizeFile return empry, csv_file:%+v", file_name)
    //     os.Exit(1)
    // }
    // result := conn.Create(bluckModels)
    // fmt.Printf("create result:%+v", result)
    
    //ParseQuizeFileToJson(file_name)
    SelectFromTableMarshalToFile(conn)
}
