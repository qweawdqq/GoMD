package controllers

import (
	"GoMD/models"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}
	// 是否开启搜索引擎查询
	isUseSearch = true
	//text  = "11111111大周日的这个点了 我还要Ò加班"
	//text1 = `111111日不知道`
	//text2 = `111111111不知道啊不知道`
	//text3 = `11111111我就是不知道这是啥`

	file, _ = exec.LookPath(os.Args[0])

	opts = types.EngineOpts{
		Using: 1,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.DocIdsIndex,
		},
		UseStore: true,
		// StoreFolder: path,
		StoreEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
		//GseDict: "../../data/dict/dictionary.txt",
		GseDict:       "testdata/test_dict.txt", //分词器
		StopTokenFile: "data/dict/stop_tokens.txt",
	}
)

func init() {
	InitEngine()
	fmt.Println("是否使用引擎进行搜索：", isUseSearch)
}

func InitEngine() {
	fmt.Println("file的路径", file)
	// var path = "./riot-index"
	//新增注册实体
	//restoreIndex()
	gob.Register(models.Article{})
	searcher.Init(opts)
	//defer searcher.Close()
	// os.MkdirAll(path, 0777)

	// Add the document to the index, docId starts at 1
	//searcher.Index(1, types.DocData{Content: text})
	//searcher.Index(2, types.DocData{Content: text1})
	//searcher.Index(3, types.DocData{Content: text2})
	//searcher.Index(5, types.DocData{Content: text3})
	////
	////searcher.RemoveDoc(5)
	////
	////// Wait for the index to refresh
	searcher.Flush()

	log.Println("Created index number: ", searcher.NumDocsIndexed())
}

func test() {
	searcher.Init(opts)
	defer searcher.Close()
}

func restoreIndex() {
	// var path = "./riot-index"
	searcher.Init(opts)
	defer searcher.Close()
	// os.MkdirAll(path, 0777)

	// Wait for the index to refresh
	searcher.Flush()

	log.Println("recover index number: ", searcher.NumDocsIndexed())
}

func SearchSomgThing(stext string) *[]models.Article {
	//initEngine()
	//test()

	fmt.Println("搜索的关键字", stext)

	ssss := types.SearchReq{
		Text: stext,
		RankOpts: &types.RankOpts{
			OutputOffset: 0,
			MaxOutputs:   150,
		}}

	sea := searcher.Search(ssss)
	fmt.Println("search response: ", sea, "; docs = ", sea.Docs)
	return requestToData(sea)
	//fmt.Println(sea.Docs)

	//os.RemoveAll("riot-index")
}

//将查询的结果转换成所需要的类型
func requestToData(output types.SearchResp) *[]models.Article {
	//var redText = "<font color=\"red\">" + stext + "</font>"
	//var redText = stext+"</br>"
	reList := []models.Article{}
	var re *models.Article
	for _, doc := range output.Docs.(types.ScoredDocs) {
		//fmt.Println("搜索时候的fmt", doc)
		//fmt.Println("content",doc.Content)//具体搜索出来的内容
		//fmt.Println("DocId",doc.DocId)
		//fmt.Println("ScoredID",doc.ScoredID)
		attri := doc.Attri.(models.Article)
		//内容里面要做出拆分
		strList := strings.Split(doc.Content, "$$")
		var text string
		if len(strList) > 1 {
			text = strList[1]
		} else {
			text = doc.Content
		}

		re = &models.Article{attri.Id, attri.Cid, attri.Title,
			attri.Author, attri.Tags, attri.Renew, text,
			attri.Summary, attri.Status,
			attri.Password, attri.Link}
		//text = strings.Replace(text, stext, redText, -1)
		//text = "<p>"+text+"</p>"
		reList = append(reList, *re)
	}
	//fmt.Println("总数",sea.NumDocs)
	return &reList
}

//标题、作者、标签
func AddContent(art *models.Article) (err error) {
	//defer searcher.Close()
	//错误
	log.Println("添加了文件标题为: ", art.Title)
	err = GetError()
	//为文件添加属性
	//attri :=types.Attri{title,author,unit.GetTimeNow(),0}

	attri := models.Article{art.Id, art.Cid, art.Title, art.Author, art.Tags, art.Renew, "", art.Summary, art.Status, art.Password, art.Link}
	content := art.Content
	arr := fmt.Sprint(attri) + "$$" + content
	ct := types.DocData{Content: arr}
	ct.Attri = attri
	//生成主键
	//docid := uint64(time.Now().Unix())
	docid := string(art.Id)
	//插入
	searcher.Index(docid, ct)
	//searcher.Index(22, types.DocData{Content: content})
	//刷新
	searcher.Flush()
	log.Println("使用了引擎添加文章功能，目前总数量为：", searcher.NumDocsIndexed())
	return err
}

func DeleteContent(id int) {
	docid := string(id)
	searcher.RemoveDoc(docid, true)
	fmt.Println("使用了引擎删除文章功能", docid)

}

func UpdateContent(id int, art *models.Article) {
	docid := string(id)
	content := art.Title + "$$" + art.Content
	ct := types.DocData{Content: content}
	ct.Attri = art
	searcher.IndexDoc(docid, ct, true)
	fmt.Println("使用了引擎更新文章功能")
}

func GetError() (err error) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic recover! p:", p)
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
			} else {
				err = errors.New("panic")
			}
			debug.PrintStack()
		}
	}()
	return err
}
func GetTimeNow() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05") //2018-7-15 15:23:00

}
